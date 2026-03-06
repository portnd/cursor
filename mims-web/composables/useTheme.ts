import {
	MenuComponent,
	ScrollComponent,
	StickyComponent,
	ToggleComponent,
	DrawerComponent,
	SwapperComponent,
} from "@/assets/themes/ts/components"

/**
 * @description Initialize KeenThemes Custom Components
 */
const initializeComponents = (): any => {
	setTimeout(() => {
		ToggleComponent.bootstrap()
		StickyComponent.bootstrap()
		MenuComponent.bootstrap()
		ScrollComponent.bootstrap()
		DrawerComponent.bootstrap()
		SwapperComponent.bootstrap()
	}, 0)
}

/**
 * @description Reinitialize KeenThemes Custom Components
 */
const reinitializeComponents = () => {
	setTimeout(() => {
		ToggleComponent.reinitialization()
		StickyComponent.reInitialization()
		MenuComponent.reinitialization()
		reinitializeScrollComponent().then(() => {
			setTimeout(() => {
				ScrollComponent.updateAll()
			}, 1000)
		})
		DrawerComponent.reinitialization()
		SwapperComponent.reinitialization()
	}, 0)
}

const reinitializeScrollComponent = (): any => {
	return Promise.resolve(ScrollComponent.reinitialization())
}

export { initializeComponents, reinitializeComponents, reinitializeScrollComponent }
