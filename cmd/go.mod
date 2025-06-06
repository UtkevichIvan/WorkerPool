module cmd

go 1.24

require (
	myModule v0.0.0
)

replace (
	myModule => ../workerpool
)