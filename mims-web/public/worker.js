self.onmessage = function (e) {
	// Perform some operation
	const result = e.data // Example operation
	registerMember(result)
}

function registerMember(conditionList) {
	const data = JSON.parse(conditionList)

	data.forEach((parent) => {
		const mappedCoordinates = parent.geom_cl.map((item) =>
			item.coordinates.map((coord) => `${coord[0]} ${coord[1]}`).join(", ")
		)

		postMessage({ geom: `MULTILINESTRING((${mappedCoordinates.join("),(")}))`, color: parent.color })
	})
}
