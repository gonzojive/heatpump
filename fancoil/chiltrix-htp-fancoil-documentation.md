# Chiltrix/HTP Fan coil documentation

CXI-35005-310586 modbus protocol

								

MODBUS RTU
From:				To:
Vesion	1.0			update:

								

## 1. Transmission Format

| Baud Rate     | 9600bps                       |
| ------------- | ----------------------------- |
| Start bit     | 1                             |
| Byte width    | 8                             |
| Parity        | N                             |
| Stop bits     | 1                             |
| Slave address | Unit's address (default is 15)|

								

## 2. Packet Format 								

| Address | Function                                                                                      | Data      | CRC checksum |
| ------- | --------------------------------------------------------------------------------------------- | --------- | ------------ |
| 16bits  | 16bits<br>03:Function of reading multi registers<br>16:Function of presenting multi registers | N\*16bits | 16bits       |
	
							

## 3. Data types								

| Data Type | Description                                                                                                                                                                                                                                                                                    |
| ---------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| TEMP1      | Signed byte ，Precision 0.1℃，Formula：T\*10，Temperature range ：\-30~97℃（if temperature is shown as 25°C, data transmission is 250 according to the preceding formula. When bit15=1 , it means minus. when bit15=0, it means integer );When this value is 32767, corresponding sensor is faulty.） |
| DIGI1      | Unsigned byte, unit is 1. When shows 123, the transmitted data is 123.                                                                                                                                                                                                                         |
| TIMER1     | Unsigned byte, unit is hours. Setting a value 1-11 starts a timer. Overwriting a timer with the same value does not reset the timer. Overwriting with a different value is not fully untested (quickly writing one value and then back to the original value does not seem to reset the timer, but it's possible delaying between writes would trigger a reset if an event loop needs to pick up the change). |

| Address | HEX  | Function Code | Content                                  | Description                                                                                                           | Remark |
| ------- | ---- | ------------- | ---------------------------------------- | --------------------------------------------------------------------------------------------------------------------- | ------ |
| 28301   | 6E8D | 03/10         | On/Off                                   | 0=off,1=on                                                                                                            | DIGI1  |
| 28302   | 6E8E | 03/10         | Mode                                     | 0～auto；1～cooling;2～dehumidification；3～ventilate；4～heating；                                                            | DIGI1  |
| 28303   | 6E8F | 03/10         | Fanspeed                                 | 2～low speed；3～medium speed； 4～high speed；<br>6～aotu                                                                   | DIGI1  |
| 28306   | 6E92 | 03/10         | Timer off                                | Timer after which to turn the unit off.                                           | TIMER1  |
| 28307   | 6E93 | 03/10         | Timer on                                 | Timer after which to turn unit on.                                               | TIMER1  |
| 28308   | 6E94 | 03/10         | Max. set temperature                     | （\-9～96）℃                                                                                                             | DIGI1  |
| 28309   | 6E95 | 03/10         | Min. set temperature                     | （\-9～96）℃                                                                                                             | DIGI1  |
| 28310   | 6E96 | 03/10         | Cooling set temperature                  |                                                                                                                       | TEMP1  |
| 28311   | 6E97 | 03/10         | heating set temperature                  |                                                                                                                       | TEMP1  |
| 28312   | 6E98 | 03/10         | Cooling set temperature at auto mode     |                                                                                                                       | TEMP1  |
| 28313   | 6E99 | 03/10         | heating set temperature at auto mode     |                                                                                                                       | TEMP1  |
| 28314   | 6E9A | 03/10         | Anti-cooling wind setting temperature    | （5～40）℃ In heating mode, if the coil temp is lower than this value, the fan motor will stop.                      | TEMP1  |
| 28315   | 6E9B | 03/10         | Whether to enable anti-hot wind function  | （1-Yes；0-No）In cooling mode, if the coil temp. is higher than 68°F, the fan motor will stop.                      | DIGI1  |
| 28316   | 6E9C | 03/10         | Whether to enable ultra-low wind function in heat mode | （1-Yes；0-No）                                                                                            | DIGI1  |
| 28317   | 6E9D | 03/10         | Whether to use vavle                     | （1-Yes；0-No）                                                                                                          | DIGI1  |
| 28318   | 6E9E | 03/10         | Whether to use floor heating             | （1-Yes；0-No）                                                                                                          | DIGI1  |
| 28319   | 6E9F | 03/10         | Whether to use Fahrenheit                | （1-℉；0-℃）                                                                                                             | DIGI1  |
| 28320   | 6EA0 | 03/10         | Master/Slave                             | （1-Yes；0-No）                                                                                                          | DIGI1  |
| 28321   | 6EA1 | 03/10         | Unit address                             | （1～99）The default value is 15                                                                                         | DIGI1  |
| 46801   | B6D1 | 04            | Room temperature                         | Only 1℃ resolution, but value is in decidegrees C (tenths of a degree)                                                                                                                       | TEMP1  |
| 46802   | B6D2 | 04            | Coil temperature                         | Only 1℃ resolution, but value is in decidegrees C (tenths of a degree)                                                                                                                       | TEMP1  |
| 46803   | B6D3 | 04            | Current fan speed                        | 0 Off；1 Ultra-low speed； 2 Low speed；3 Medium speed；4 High speed；5 Top speed；6 Auto                                   | DIGI1  |
| 46804   | B6D4 | 04            | Fan revolution                           | 0～2000 （rpm）                                                                                                          | DIGI1  |
| 46805   | B6D5 | 04            | Electromagnetic Valve                    | 0 Off； 1 On                                                                                                           | DIGI1  |
| 46806   | B6D6 | 04            | Remote on/off                            | 0 Open;1 close                                                                                                        | DIGI1  |
| 46807   | B6D7 | 04            | Simulation signal                        | 0 (The main engine needs to be switched to non-heating mode)；1 (The main engine needs to be switched to heating mode) | DIGI1  |
| 46808   | B6D8 | 04            | Fan speed signal feedback fault          | （1-Yes；0-No）                                                                                                          | DIGI1  |
| 46809   | B6D9 | 04            | Room temperature sensor fault            | （1-Yes；0-No）                                                                                                          | DIGI1  |
| 46810   | B6DA | 04            | Coil temperature sensor fault            | （1-Yes；0-No）                                                                                                          | DIGI1  |


**Note：We can't set P27 as 2 if P17 is set as 0(Only the unit with valve can be set)**
