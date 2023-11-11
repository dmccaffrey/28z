package core

var Variables = map[string]CoreValue{
	// Math
	"g":     FloatValue{value: 9.80665},
	"tau":   FloatValue{value: 6.2831855},
	"pi":    FloatValue{value: 3.1415926},
	"phi":   FloatValue{value: 1.6180339},
	"e":     FloatValue{value: 2.7182818},
	"gauss": FloatValue{value: 0.8346268},
	"c":     FloatValue{value: 299792458},

	// Inernal
	"ram-bytes":     FloatValue{value: 8192},
	"render-width":  FloatValue{value: 92},
	"render-height": FloatValue{value: 30},
	"render-bytes":  FloatValue{value: 2760},
}
