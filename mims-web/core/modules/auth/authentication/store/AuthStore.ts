import jwtDecode from "jwt-decode"
import { IUser, RefreshTokenService } from "../infrastructure"

export interface IInverval {
	now: string | null
	exp: string | null
	minutes: number
}

export interface IRefreshToken {
	exp: number
	refresh_uuid: string
	token: string
	user_id: number
}

export interface IState {
	accessToken: string
	refreshToken: string
	user: IUser | null
	inverval: IInverval | null
}

export const useAuthStore = defineStore("auth", {
	state: (): IState => ({
		accessToken: "",
		refreshToken: "",
		user: null,
		inverval: null,
	}),
	actions: {
		autoRefresh() {
			const config = useRuntimeConfig()

			// เรียกใช้ตอน Refresh หน้าเว็บไซต์
			setTimeout(() => {
				this.refresh()
			}, 10 * 1000)

			// เรียกใช้แบบ Auto
			console.log("[JWT] Interval (" + Number(config.autoRefreshJwtInterval) + "min)")
			setInterval(() => {
				this.refresh()
			}, Number(config.autoRefreshJwtInterval) * 60 * 1000)

			return Promise.resolve()
		},
		async refresh() {
			try {
				if (this.isLoggedIn === true) {
					const config = useRuntimeConfig()

					// Calculate expiration time
					const now = new Date(Date.now())
					const exp = new Date((this.user?.exp || Date.now()) * 1000)
					const minutes = (exp.getTime() - now.getTime()) / 1000 / 60

					// เก็บข้อมูลเวลา
					this.inverval = {
						now: String(now),
						exp: String(exp),
						minutes,
					}

					console.log("[JWT] Calculate expiration time... ✅")
					console.log(`[JWT] Expiration time (${toNumber(minutes, 2)}/${Number(config.jwtExpirationTime)} min)`)

					// คำนวณเวลาหมดอายุของ AccessToken
					if (minutes < Number(config.jwtExpirationTime)) {
						const refreshTokenService = new RefreshTokenService()
						const res = await refreshTokenService.refreshToken({
							access_token: this.accessToken,
							refresh_token: this.refreshToken,
						})

						if (res.status !== false) {
							// กรณี Refresh token สำเร็จ
							console.log("[JWT] Refreshing a Token...")
							this.setToken(res.data.access_token, res.data.refresh_token)
						} else {
							console.log("[JWT] Failed token refresh 🔥")
							this.logout()
						}
					}
				}
			} catch (error) {
				console.log("[JWT] Refresh token 🔥🔥", error)
				this.logout()
			}
		},
		setToken(accessToken: string, refreshToken: string) {
			this.accessToken = accessToken
			this.refreshToken = refreshToken

			// เก็บข้อมูลของ User
			if (accessToken !== "") {
				// Decode ข้อมูลจาก AccessToken
				const decodedAccessToken: IUser = jwtDecode(accessToken)
				this.user = decodedAccessToken

				console.log("[User] Store updated ✅")
			}
		},
		logout() {
			this.accessToken = ""
			this.refreshToken = ""
			this.user = null

			return navigateTo("/auth/login")
		},
		subscribeStore() {
			const self = this
			self.$subscribe(() => {
				const isLoggedIn = () => {
					return this.isLoggedIn
				}

				// กรณีเข้าสู่ระบบแล้ว
				if (!isLoggedIn()) {
					return useAuthStore().logout()
				}

				// กรณียังไม่ได้เข้าสู่ระบบ
				watch(isLoggedIn, (value) => {
					if (value === false) {
						return useAuthStore().logout()
					}
				})
			})
		},
	},
	getters: {
		isLoggedIn(state): boolean {
			const date = new Date()
			const now = Math.floor(date.getTime() / 1000)

			// เวลาหมดอายุ
			if (state.user?.exp) {
				return state.accessToken !== "" && state.refreshToken !== "" && state.user?.exp > now
			}

			return false
		},
	},
	persist: true,
})
