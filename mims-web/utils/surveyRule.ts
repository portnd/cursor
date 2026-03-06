export const getGrade = (id: number) => {
	const refGrade = useInitData().refGrade?.()
	const { color = "", name = "" } = refGrade?.find((el: any) => el.id === id) ?? {}

	return { color, name }
}
