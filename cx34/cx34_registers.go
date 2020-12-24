package cx34

// This file is used to assign names to modbus registers.

/*

Table of registers with values that changed

| Register no. | Notes                       | Value |
|--------------|-----------------------------|-------|
| 143          | AC heating setpoint         | 39    |
| 144          | Domestic Hot Water setpoint | 51    |
| 200          | C0                          | 22    |
| 201          | C1                          | 700   |
| 202          | C2                          | 49    |
| 203          | C3                          | 8     |
| 204          | C4                          | 451   |
| 205          | C5                          | 453   |
| 206          | C6                          | 249   |
| 207          | C7                          | 245   |
| 208          | C8                          | 247   |
| 209          | C9                          | 60    |
| 213          | C13                         | 110   |
| 227          | C27                         | 43    |
| 229          | C29                         | 1     |
| 235          | C35                         | 0     |
| 237          | C37                         | 859   |
| 240          | C40                         | 98    |
| 241          | C41                         | 237   |
| 243          | C43                         | 0     |
| 245          | C45                         | 895   |
| 248          | C48                         | 6     |
| 250          | C50                         | 55    |
| 251          | C51                         | 15    |
| 255          | C55                         | 239   |
| 256          | C56                         | 55    |
| 257          | C57                         | 69    |
| 258          | C58                         | 389   |
| 260          | C60                         | 24    |
| 261          | C61                         | 14    |
| 280          | C80                         | 248   |
| 281          | C81                         | 399   |
| 282          | C82                         | 399   |
*/

// Known Register values.
const (
	TargetACHeatingModeTemp    Register = 143
	TargetDomesticHotWaterTemp Register = 144
	// See page 51 of https://www.chiltrix.com/documents/CX34-IOM-3.pdf
	OutPipeTemp               Register = 200
	CompressorDischargeTemp   Register = 201
	AmbientTemp               Register = 202
	SuctionTemp               Register = 203
	PlateHeatExchangerTemp    Register = 204
	ACOutletWaterTemp         Register = 205
	SolarTemp                 Register = 206
	CompressorCurrentValueP15 Register = 209 // 0.00-30.0A
	UsageSideWaterFlowRate    Register = 213 // tenths of a liter per minute
	P03Status                 Register = 214
	CompressorFrequency       Register = 227

	WaterOutletSensorTemp1 Register = 204
	WaterOutletSensorTemp2 Register = 205
	WaterFlowRate          Register = 213
	WaterInletSensorTemp1  Register = 281
	WaterInletSensorTemp2  Register = 282
)

// Source: https://www.chiltrix.com/control-options/Remote-Gateway-BACnet-Guide-rev2.pdf
var registerNames = map[Register]string{
	1:   "Power-Down Recovery Function",
	2:   "Single / Three Phase Selection",
	3:   "Power Frequency",
	4:   "Heat Source Selection",
	5:   "Heating Temp Control Method",
	6:   "Defrost Method Selections",
	7:   "Freecooling Validation",
	8:   "Frequency Control Method",
	9:   "DHW Validation",
	10:  "Air Cond And Heating Validation",
	11:  "Air Cond And Cooling Validation",
	12:  "DHW Hot Water Temp Hysteresis",
	13:  "AC Temp Hysteresis",
	14:  "Fan Motor Category",
	15:  "Maximum Speed Of The Fan",
	16:  "Heating Fan Speed Control Temp Diff",
	17:  "Cooling Fan Speed Control Temp Diff",
	18:  "Defrost Method",
	19:  "Defrost Starting Temp",
	20:  "Defrost Interval Time Multiple Rate",
	21:  "The First Defrost Interval",
	22:  "Defrost Exist Temp",
	23:  "Hot Water Frequency Limitation",
	24:  "AC Heating Au Mode Highest Temp",
	25:  "AC Heating Au Mode Offset Temp",
	26:  "Solenoid Valve Function Parameters",
	27:  "C4 Water Pump Type Selection",
	28:  "Water Pump Working Mode",
	29:  "EC Water Pump C4 Minimum Speed",
	30:  "C5 Water Pump Type Selection",
	31:  "DHW E-Heater Activated Ambient Temp",
	32:  "Electric Heating Function",
	33:  "AC E-Heater Activated Ambient Temp",
	34:  "2nd Heat Source Starting Air Temp",
	35:  "AC Anti-Freezing Temp",
	36:  "Virus Killing Interval Days",
	37:  "Start Virus Killing Time",
	38:  "Virus Killing Holding Time",
	39:  "Target Temp Of Virus Killing",
	40:  "AC Water Flow Switch Type Selection",
	41:  "AC Minimum Water Flow",
	42:  "Water Src Water Flow Switch Type Sel",
	43:  "Lowest Water Flow Of Water Source",
	44:  "Air Src Heat Pump Freecooling Func",
	45:  "Air Src Freecooling Function",
	46:  "Cooling Maximum Set Temp",
	47:  "Heating Maximum Set Temp",
	48:  "DHW The Highest Set Temp",
	49:  "Debugging Fixed Operating Frequency",
	50:  "Run Setting Frequency",
	51:  "EEV Manually Open Degree (Heating)",
	52:  "EEV Manually Open Degree (Cooling)",
	53:  "EEV Control Mode",
	54:  "Target Overheat Degree (Heating)",
	55:  "Target Overheat Degree (Cooling)",
	56:  "Night Mode Validation",
	57:  "Night Mode Starting Point",
	58:  "Night Mode Ending Point",
	59:  "Model Selection",
	60:  "Use High And Low Pressure Transmitter",
	61:  "Temp Diff To Ctrl C4 Water Pump Speed",
	62:  "Compressor Manufacturer",
	63:  "Forced Sterilization",
	64:  "System Parameter Recovery",
	65:  "Compressor Manufacturer 2",
	66:  "Virus Killing Function Validation",
	67:  "EEV Max Manual Open",
	68:  "Defrosting EEV Manual Open",
	69:  "AC Electric Heater Power W",
	70:  "C Or F Degree",
	71:  "Heat Recovery Function Validation",
	72:  "AC Rated Voltage",
	73:  "AC Heat Transfer Coefficient",
	74:  "AC Voltage Compensation",
	75:  "Cooling Inlet Target Temp Range",
	76:  "AC Heating Minimum Frequency",
	77:  "Own 485 Address",
	78:  "Error Recovery",
	79:  "Switch On/Off",
	80:  "Operating Mode",
	81:  "AC Cooling Target Temp",
	82:  "AC heating Target Temp",
	83:  "Hot Water Target Temp",
	84:  "AC Heating Au Mode",
	85:  "Hot Water Au Mode",
	86:  "Out Pipe Temp",
	87:  "Compressor Discharge Temp",
	88:  "Ambient Temp",
	89:  "Suction Temp",
	90:  "Plate Heat Exchanger Inlet Temp",
	91:  "AC Outlet Water Temp",
	92:  "Solar Temp",
	93:  "Compressor Current Value",
	94:  "Usage Side Water Flow Volume",
	95:  "P03 Status",
	96:  "P04 Status",
	97:  "P05 Status",
	98:  "P06 Status",
	99:  "P07 Status",
	100: "P08 Status",
	101: "P09 Status",
	102: "P10 Status",
	103: "High Pressure Switch Status",
	104: "Low Pressure Switch Status",
	105: "Second High Pressure Switch Status",
	106: "Inner Water Flow Switch",
	107: "Compressor Frequency",
	108: "Overheat Switch Status",
	109: "Outdoor Fan Motor",
	110: "Electrical Valve 1",
	111: "Electrical Valve 2",
	112: "Electrical Valve 3",
	113: "Electrical Valve 4",
	114: "C4Water Pump",
	115: "C5Water Pump",
	116: "C6Water Pump",
	117: "Accum Days After Last Virus Killing",
	118: "Outdoor Modular Temp",
	119: "Expansion Valve 1 Opening Degree",
	120: "Expansion Valve 2 Opening Degree",
	121: "Inner Pipe Temp Display",
	122: "Heating Method 2 Target Temp",
	123: "Run Returning Lubrication Oil Func",
	124: "Fan Type",
	125: "EC Fan Motor 1 Speed",
	126: "EC Fan Motor 2 Speed",
	127: "Water Pump Types",
	128: "Water Pump1 Speed",
	129: "Water Pump2 Speed",
	130: "Inductor AC Current Value",
	131: "Driver Working Status Value",
	132: "Compressor Shut Down Code",
	133: "Driver Allowed Highest Frequency",
	134: "Reduce Frequency Temp Setting",
	135: "Input AC Voltage Value",
	136: "Input AC Current Value",
	137: "Compressor Phase Current Value",
	138: "Bus Line Voltage",
	139: "Fan Shutdown Code",
	140: "Ipm Temp",
	141: "Compressor Total Running Time",
	142: "E-Heater Compensation Power",
	143: "TargetACHeatingModeTemp",
	144: "TargetDomesticHotWaterTemp", // was: "Din7 AC Cooling Mode Switch",
	145: "DHW Current Temp",
	146: "AC Heating Current Temp",
	147: "AC Cooling Current Temp",
	148: "Error Unit1 Err1",
	149: "Error Unit2 Err2",
	150: "Error Unit3 Err3",
	151: "Error Unit4 Err4",
	152: "Error Unit5 Err5",
	153: "Error Unit5 Err6",
	154: "Comp Discharge High Temp Protection",
	155: "Outdoor Air Temp Sen Error",
	156: "Outer Coil Pipe Temp Sen Error",
	157: "Pipe Returned Gas Sen Error",
	158: "Indoor Refrigerant Pipe Temp Sen Err",
	159: "Coil High Temp Protection",
	160: "Solar Water Temp Sen Error",
	161: "AC Inlet Water Temp Sen Error",
	162: "AC Outlet Water Temp Sen Error",
	163: "DHW Temp Sen Error",
	164: "Indoor Ambient Sen Error",
	165: "Water Src Inlet Water Temp Sen Error",
	166: "Water Src Outlet Temp Sen",
	167: "System Anti Freeze Twice",
	168: "DHW Anti Freeze Twice",
	169: "Discharge Probe Error",
	170: "High Pressure Protection",
	171: "Low Pressure Protection",
	172: "Comp Overheat Protection",
	173: "Over Current Protection",
	174: "Indoor Unit Water Flow Error",
	175: "Outdoor Water Flow Error",
	176: "Miss Phase",
	177: "Wrong Phase",
	178: "Com Error",
	179: "Water Src Anti Freeze",
	180: "Water Src Water Flow Not Enough",
	181: "Voltage Protection",
	182: "Ipm Fault",
	183: "Comp Drive Fault",
	184: "Comp Over Current Protection 1",
	185: "ERR3.0",
	186: "Ipm Overheat",
	187: "PFC Fault",
	188: "DC Bus Overvoltage",
	189: "DC Bus Undervoltage",
	190: "AC Input Over Or Under Voltage",
	191: "AC Input Current Protection",
	192: "Temperature Sen Fault",
	193: "DSO And Mainboard Com Fault",
	194: "Control Board And Inverter Com Fault",
	195: "Inlet/Outlet Wtr Temp Diff Is Too Big",
	196: "AC System Antifreeze Twice",
	197: "ERR3.12",
	198: "ERR3.13",
	199: "Ctrl Panel Param Are Not Initialized",
	// Starting at 200, it's all the C parameters from the details screen.
	200: "ERR3.15",
	201: "EC Fan 1 Fault",
	202: "EC Fsn 2 Fault",
	203: "Heat Recovery Warning",
	204: "WaterOutletSensorTemp1",
	205: "WaterOutletSensorTemp2",
	213: "WaterFlowRate", // tenths of a liter / minute
	281: "WaterInletSensorTemp1",
	282: "WaterInletSensorTemp2",
}
