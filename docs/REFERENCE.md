# Reference

## Loxone Control Types

### `AalEmergency`
Emergency alarm control.

- **States:**
  - `status` (0: running, 1: alarm triggered, 2: reset input asserted, 3: app disabled)
  - `disableEndTime` (unix timestamp)
  - `resetActive` (text state)
- **Commands:**
  - `trigger`: Trigger an alarm from the app
  - `quit`: Quit an active alarm
  - `disable/<timespan>`: Disable the control for the given time in seconds.

### `AalSmartAlarm`
Smart Alarm control.

- **States:**
  - `alarmLevel` (0=no alarm, 1=immediate, 2=delayed)
  - `alarmCause` (String)
  - `isLocked`
  - `isLeaveActive`
  - `disableEndTime`
- **Commands:**
  - `confirm`: Confirm pending alarm
  - `disable/{seconds}`: Disable control for a certain period
  - `startDrill`: Execute test alarm

### `Alarm`
Burglar alarm.

- **States:**
  - `armed`
  - `nextLevel` (ID of next alarm level)
  - `nextLevelAt` (unix timestamp)
  - `nextLevelDelayTotal`
  - `disabledMove`
  - `startTime`
- **Commands:**
  - `on`: Arms the AlarmControl
  - `on/{number}`: 0=arm without movement, 1=arm with movement
  - `delayedon`: Arms with delay
  - `delayedon/{number}`: 0=delayed arm without movement, 1=delayed arm with movement
  - `off`: Disarms the AlarmControl
  - `quit`: Acknowledge the alarm
  - `dismv/{number}`: 0=disable movement, 1=enable movement

### `AlarmChain`
Alarm Sequence.

- **States:**
  - `activeAlarmType` (Bitmap: 0: Inactive, 1: Acknowledged, 2: Alarm, 4: Urgent, 8: EMS)
  - `nextAlarmLevelAt`
  - `activeAlarmText`
  - `nextAlarmText`
  - `iterationCount`
- **Commands:**
  - `quit`: Acknowledge the alarm

### `AlarmClock`
Alarm clock control.

- **States:**
  - `isEnabled`
  - `isAlarmActive`
  - `confirmationNeeded`
  - `entryList` (JSON Object)
  - `currentEntry`
  - `nextEntry`
  - `nextEntryMode`
  - `ringingTime`
  - `ringDuration`
  - `snoozeTime`
  - `snoozeDuration`
  - `nextEntryTime`
- **Commands:**
  - `snooze`: Snoozes the current active entry
  - `dismiss`: Dismisses the current active entry
  - `setSnoozeDuration/{number}`: Sets the snoozing duration

### `AudioZone`
Music Server Zone.

- **States:**
  - `serverState` (-3 to 2)
  - `playState` (-1=unknown, 0=stopped, 1=paused, 2=playing)
  - `clientState` (0=offline, 1=initializing, 2=online)
  - `power`
  - `volume`
  - `maxVolume`
  - `shuffle`
  - `repeat`
  - `songName`
  - `duration`
  - `progress`
  - `album`
  - `artist`
  - `station`
  - `genre`
  - `cover`
  - `source`
- **Commands:**
  - `volume/{newVolume}`
  - `volstep/{step}`
  - `prev`
  - `next`
  - `play`
  - `pause`
  - `shuffle`
  - `repeat/{repeatState}`
  - `on`
  - `off`
  - `source/{sourceNumber}`

### `AudioZoneV2`
Music Server Zone Gen 2.

- **States:**
  - `serverState`
  - `playState`
  - `clientState`
  - `power`
  - `volume`
  - `isLocked`
- **Commands:**
  - `volUp`
  - `volDown`
  - `volume/{value}`
  - `prev`
  - `next`
  - `play`
  - `Pause`

### `CarCharger`
Wallbox Gen 1.

- **States:**
  - `status` (0=Offline, 1=Initializing, 2=Online)
  - `charging` (0,1)
  - `connected` (0,1)
  - `chargingFinished` (0,1)
  - `power` (kW)
  - `energySession` (kWh)
  - `limitMode` (0=Off, 1=Manual, 2=Automatic)
  - `currentLimit` (kW)
  - `chargeDuration` (Secs)
- **Commands:**
  - `charge/on`: Start charging
  - `charge/off`: Stop/pause charging
  - `limitMode/{mode}`
  - `limit/{limit}`

### `ClimateController`
Climate Controller.

- **States:**
  - `currentMode` (0=No req, 1=Heat, 2=Cool, 3=HeatBoost, 4=CoolBoost, 5=Service, 6=ExtHeater)
  - `autoMode` (-1=Off, 0=Heat+Cool, 1=Heat, 2=Cool)
  - `currentAutomatic`
  - `heatingTempBoundary`
  - `coolingTempBoundary`
  - `actualOutdoorTemp`
  - `averageOutdoorTemp`
  - `serviceMode`
  - `ventilation`
  - `excessEnergy`
- **Commands:**
  - `setServiceMode/{active}`
  - `ventilation/{active}`
  - `autoMode/{mode}`
  - `setHeatingBoundary/{temp}`
  - `setCoolingBoundary/{temp}`

### `ClimateControllerUS`
HVAC Controller.

- **States:**
  - `mode` (0=Off, 1=Heat+Cool, 2=Heat, 3=Cool)
  - `currentStatus` (Bitmask)
  - `fan`
  - `humidity`
  - `actualOutdoorTemp`
  - `minimumTempCooling`
  - `maximumTempHeating`
  - `stage`
  - `serviceMode`
- **Commands:**
  - `ventilation/{active}`
  - `setMode/{mode}`
  - `useEmergency/{1/0}`
  - `setMinimumTempCooling/{temp}`
  - `setMaximumTempHeating/{temp}`
  - `setServiceMode/{active}`

### `ColorPickerV2`
Color Picker for LightControllerV2.

- **States:**
  - `color` (hsv or temp string)
  - `sequence`
  - `sequenceColorIdx`
- **Commands:**
  - `setFav/{favColorIdx}/{color}`
  - `hsv(hue, sat, val)`
  - `temp(brightness, temperature)`
  - `setBrightness/{value}`

### `Daytimer`
Schedule.

- **States:**
  - `mode`
  - `override`
  - `value`
  - `entriesAndDefaultValue`
  - `resetActive`
- **Commands:**
  - `pulse`
  - `default`
  - `startOverride/{value}/{howLongInSecs}`
  - `stopOverride`
  - `set` (Set entries)

### `Dimmer`
Dimmer.

- **States:**
  - `position`
  - `min`
  - `max`
  - `step`
- **Commands:**
  - `on`
  - `off`
  - `{pos}`: Set position

### `EnergyFlowMonitor`
Energy Flow Monitor.

- **States:**
  - `Ppwr`: Production Power
  - `Gpwr`: Grid Power
  - `Spwr`: Storage Power
  - `Pre`: Price export
  - `Pri`: Price import
  - `CO2`: CO2 Factor
  - `actual0`...`actualN`
- **Commands:**
  - `get/{viewType}`

### `Gate`
Gate control.

- **States:**
  - `position` (0=closed, 1=open)
  - `active` (-1=close, 0=not moving, 1=open)
  - `preventOpen`
  - `preventClose`
- **Commands:**
  - `open`
  - `close`
  - `stop`
  - `partiallyOpen`

### `Hourcounter`
Maintenance counter.

- **States:**
  - `total`
  - `remaining`
  - `lastActivation`
  - `overdue`
  - `maintenanceInterval`
  - `active`
- **Commands:**
  - `reset`: Reset remaining/overdue
  - `resetAll`: Reset all stats

### `InfoOnlyAnalog`
Virtual state (Analog).

- **States:**
  - `value`

### `InfoOnlyDigital`
Virtual state (Digital).

- **States:**
  - `value` (0 or 1)
  - `text` (on/off text)
  - `color` (on/off color)

### `InfoOnlyText`
Virtual state (Text).

- **States:**
  - `text`

### `Intelligent Room Controller v2`
Intelligent Room Controller V2.

- **States:**
  - `activeMode` (0=Eco, 1=Comfort, 2=BuildingProt, 3=Manual, 4=Off)
  - `operatingMode` (0=Auto H+C, 1=Auto H, 2=Auto C, 3=Manual H+C, 4=Manual H, 5=Manual C)
  - `tempActual`
  - `tempTarget`
  - `comfortTemperature`
  - `comfortTemperatureCool`
  - `comfortTolerance`
  - `openWindow`
  - `co2`
  - `humidityActual`
- **Commands:**
  - `override/{modeId}/{until}/{temp}`
  - `stopOverride`
  - `setComfortTemperature/{temp}`
  - `setManualTemperature/{temp}`
  - `setOperatingMode/{opMode}`

### `Intercom`
Door Controller.

- **States:**
  - `bell` (0=not ringing, 1=ringing)
  - `lastBellEvents`
  - `version`
- **Commands:**
  - `answer`

### `IntercomV2`
Intercom V2.

- **States:**
  - `bell`
  - `missedBellEvents`
  - `deviceState`
- **Commands:**
  - `answer`
  - `mute/{0/1}`

### `Jalousie`
Blinds/Shading.

- **States:**
  - `up`
  - `down`
  - `position` (0=upper, 1=lower)
  - `shadePosition` (0=horizontal, 1=vertical)
  - `safetyActive`
  - `autoAllowed`
  - `autoActive`
  - `locked`
- **Commands:**
  - `up`
  - `down`
  - `FullUp`
  - `FullDown`
  - `shade`
  - `auto`
  - `NoAuto`
  - `manualPosition/{targetPosition}`
  - `manualLamelle/{targetPosition}`
  - `stop`

### `LightControllerV2`
Lighting Controller V2.

- **States:**
  - `activeMoods` (List of active mood IDs)
  - `moodList` (List of available moods)
  - `favoriteMoods`
  - `additionalMoods`
  - `presence`
- **Commands:**
  - `changeTo/{moodId}`
  - `addMood/{moodId}`
  - `removeMood/{moodId}`
  - `plus`: Next mood
  - `minus`: Previous mood
  - `presence/{on,off}`

### `Meter`
Utility meter.

- **States:**
  - `actual`
  - `total`
  - `totalNeg` (if bidirectional)
- **Commands:**
  - `reset`

### `NFC Code Touch`
NFC Code Touch.

- **States:**
  - `codeDate`
  - `deviceState`
  - `nfcLearnResult`
- **Commands:**
  - `code/create/...`
  - `code/update/...`
  - `code/delete/{uuid}`
  - `nfc/startlearn`
  - `nfc/stoplearn`

### `PresenceDetector`
Presence Detector.

- **States:**
  - `active`
  - `locked`
  - `activeSince`
- **Commands:**
  - `{value}`: Set value for virtual input
  - `time/{value}`: Set overrun time

### `Pushbutton`
Push button.

- **States:**
  - `active`
- **Commands:**
  - `on`
  - `pulse`

### `Radio`
Radio buttons.

- **States:**
  - `activeOutput` (ID of active output, 0=none)
- **Commands:**
  - `reset`: Deselects current
  - `{ID}`: Selects output ID
  - `next`
  - `prev`

### `Remote`
Media controller.

- **States:**
  - `mode` (0=no mode, >0 mode key)
  - `active`
- **Commands:**
  - `mode/{modeID}`
  - `on`
  - `off`
  - `play`, `pause`, `stop`, `volplus`, `volminus`, etc.

### `Sauna`
Sauna controller.

- **States:**
  - `active`
  - `power`
  - `tempActual`
  - `tempTarget`
  - `fan`
  - `drying`
  - `error`
  - `saunaError`
- **Commands:**
  - `on`
  - `off`
  - `fanon`
  - `fanoff`
  - `temp/{target}`

### `SmokeAlarm`
Fire/Water alarm.

- **States:**
  - `nextLevel`
  - `level`
  - `acousticAlarm`
  - `testAlarm`
  - `alarmCause`
- **Commands:**
  - `mute`
  - `confirm`
  - `startDrill`

### `SolarPumpController`
Thermal Solar Controller.

- **States:**
  - `bufferTemp{n}`
  - `bufferState{n}`
  - `collectorTemp`

### `Switch`
Basic relay or toggle switch.

- **States:**
  - `active` (1=On, 0=Off)
  - `jLocked`
  - `lockedOn`
- **Commands:**
  -  `On`: Turns the switch permanently On
  -  `Off`: Turns the switch permanently Off
  -  `Pulse`: Toggles the state (simulates a button press).

### `TextState`
Status text.

- **States:**
  - `textAndIcon`

### `TextInput`
Virtual text input.

- **States:**
  - `text`
- **Commands:**
  - `{text}`: Set new value

### `TimedSwitch`
Stairwell/Comfort switch.

- **States:**
  - `deactivationDelayTotal`
  - `deactivationDelay`
- **Commands:**
  - `on`
  - `off`
  - `pulse`

### `Tracker`
Tracker.

- **States:**
  - `entries` (Pipe separated string)

### `UpDownLeftRight digital`
Virtual Input (Directional).

- **Commands:**
  - `UpOn`, `UpOff`, `PulseUp`
  - `DownOn`, `DownOff`, `PulseDown`

### `ValueSelector`
Push-button +/-.

- **States:**
  - `value`
  - `min`
  - `max`
  - `step`
- **Commands:**
  - `{value}`

### `Ventilation`
Ventilation.

- **States:**
  - `speed`
  - `mode`
  - `temperatureIndoor`
  - `temperatureOutdoor`
  - `temperatureTarget`
  - `airQualityIndoor`
  - `humidityIndoor`
  - `presence`
- **Commands:**
  - `setTimer/{interval}/{speed}/{modeId}/{timerProfileIdx}`
  - `setTimer/0`: Stop timer
  - `setPresenceMin/{value}`
  - `setPresenceMax/{value}`

### `Wallbox2`
Wallbox Gen 2.

- **States:**
  - `connected`
  - `enabled`
  - `active`
  - `limit`
  - `mode`
  - `session` (JSON)
  - `pricePerkWh`
- **Commands:**
  - `allow/on`
  - `allow/off`
  - `setmode/{value}` (1-5, 99=manual)
  - `manualLimit/{value}`

### `Window`
Window.

- **States:**
  - `position` (0..1)
  - `direction` (-1=closing, 0=still, 1=opening)
  - `lockedReason`
- **Commands:**
  - `open/on`, `open/off`
  - `close/on`, `close/off`
  - `fullopen`
  - `fullclose`
  - `moveToPosition/{pos}`
  - `stop`

### `WindowMonitor`
Window and Door Monitor.

- **States:**
  - `windowStates` (List of states)
  - `numOpen`
  - `numClosed`
  - `numTilted`
  - `numLocked`
  - `numUnlocked`

## MQTT Topics
Check topic structure in [Architecture](docs/ARCHITECTURE.md).

## `<control-type>_<state>` Topics

These topics represent the current state of a specific control type.
The topic name is constructed as: `loxone/<serial>/<room>/<control>/<type>_<state>`.
Below is the reference for the last part: `<type>_<state>`.

### Common Payload Format
All state topics use the following JSON structure:
```json
{
    "value": <value>,          // The value of the state (type depends on state)
    "ts": "2024-10-01T12:34:56Z" // ISO 8601 timestamp of when the event was received
}
```

### `AalEmergency`
- **`aalemergency_status`**: Number (0-3)
- **`aalemergency_disableendtime`**: Number (Unix Timestamp)
- **`aalemergency_resetactive`**: String

### `AalSmartAlarm`
- **`aalsmartalarm_alarmlevel`**: Number (0-2)
- **`aalsmartalarm_alarmcause`**: String
- **`aalsmartalarm_islocked`**: Boolean
- **`aalsmartalarm_isleaveactive`**: Boolean
- **`aalsmartalarm_disableendtime`**: Number (Unix Timestamp)

### `Alarm`
- **`alarm_armed`**: Boolean
- **`alarm_nextlevel`**: Number (ID)
- **`alarm_nextlevelat`**: Number (Unix Timestamp)
- **`alarm_nextleveldelaytotal`**: Number (Seconds)
- **`alarm_disabledmove`**: Boolean
- **`alarm_starttime`**: Number (Unix Timestamp)

### `AlarmChain`
- **`alarmchain_activealarmtype`**: Number (Bitmap)
- **`alarmchain_nextalarmlevelat`**: Number (Unix Timestamp)
- **`alarmchain_activealarmtext`**: String
- **`alarmchain_nextalarmtext`**: String
- **`alarmchain_iterationcount`**: Number

### `AlarmClock`
- **`alarmclock_isenabled`**: Boolean
- **`alarmclock_isalarmactive`**: Boolean
- **`alarmclock_confirmationneeded`**: Boolean
- **`alarmclock_entrylist`**: JSON Object
- **`alarmclock_currententry`**: Number (ID)
- **`alarmclock_nextentry`**: Number (ID)
- **`alarmclock_nextentrymode`**: Number
- **`alarmclock_ringingtime`**: Number (Seconds)
- **`alarmclock_ringduration`**: Number (Seconds)
- **`alarmclock_snoozetime`**: Number (Seconds)
- **`alarmclock_snoozeduration`**: Number (Seconds)
- **`alarmclock_nextentrytime`**: Number (Unix Timestamp)

### `AudioZone`
- **`audiozone_serverstate`**: Number (-3 to 2)
- **`audiozone_playstate`**: Number (-1 to 2)
- **`audiozone_clientstate`**: Number (0-2)
- **`audiozone_power`**: Boolean
- **`audiozone_volume`**: Number
- **`audiozone_maxvolume`**: Number
- **`audiozone_shuffle`**: Boolean
- **`audiozone_repeat`**: Number
- **`audiozone_songname`**: String
- **`audiozone_duration`**: Number
- **`audiozone_progress`**: Number
- **`audiozone_album`**: String
- **`audiozone_artist`**: String
- **`audiozone_station`**: String
- **`audiozone_genre`**: String
- **`audiozone_cover`**: String (URL path)
- **`audiozone_source`**: Number

### `AudioZoneV2`
- **`audiozonev2_serverstate`**: Number
- **`audiozonev2_playstate`**: Number
- **`audiozonev2_clientstate`**: Number
- **`audiozonev2_power`**: Boolean
- **`audiozonev2_volume`**: Number
- **`audiozonev2_islocked`**: Boolean

### `CarCharger`
- **`carcharger_status`**: Number (0-2)
- **`carcharger_charging`**: Boolean
- **`carcharger_connected`**: Boolean
- **`carcharger_chargingfinished`**: Boolean
- **`carcharger_power`**: Number (kW)
- **`carcharger_energysession`**: Number (kWh)
- **`carcharger_limitmode`**: Number (0-2)
- **`carcharger_currentlimit`**: Number (kW)
- **`carcharger_chargeduration`**: Number (Seconds)

### `ClimateController`
- **`climatecontroller_currentmode`**: Number (0-6)
- **`climatecontroller_automode`**: Number (-1 to 2)
- **`climatecontroller_currentautomatic`**: Number
- **`climatecontroller_heatingtempboundary`**: Number
- **`climatecontroller_coolingtempboundary`**: Number
- **`climatecontroller_actualoutdoortemp`**: Number
- **`climatecontroller_averageoutdoortemp`**: Number
- **`climatecontroller_servicemode`**: Number
- **`climatecontroller_ventilation`**: Number (State)
- **`climatecontroller_excessenergy`**: Number (Bitmask)

### `ClimateControllerUS`
- **`climatecontrollerus_mode`**: Number (0-3)
- **`climatecontrollerus_currentstatus`**: Number (Bitmask)
- **`climatecontrollerus_fan`**: Number
- **`climatecontrollerus_humidity`**: Number
- **`climatecontrollerus_actualoutdoortemp`**: Number
- **`climatecontrollerus_minimumtempcooling`**: Number
- **`climatecontrollerus_maximumtempheating`**: Number
- **`climatecontrollerus_stage`**: Number
- **`climatecontrollerus_servicemode`**: Number

### `ColorPickerV2`
- **`colorpickerv2_color`**: String (e.g., "hsv(0,100,100)")
- **`colorpickerv2_sequence`**: JSON Object
- **`colorpickerv2_sequencecoloridx`**: Number

### `Daytimer`
- **`daytimer_mode`**: Number
- **`daytimer_override`**: Number (Remaining time)
- **`daytimer_value`**: Number (Analog value or 0/1)
- **`daytimer_entriesanddefaultvalue`**: JSON Object
- **`daytimer_resetactive`**: Boolean

### `Dimmer`
- **`dimmer_position`**: Number (0.0 - 1.0)
- **`dimmer_min`**: Number
- **`dimmer_max`**: Number
- **`dimmer_step`**: Number

### `EnergyFlowMonitor`
- **`energyflowmonitor_ppwr`**: Number
- **`energyflowmonitor_gpwr`**: Number
- **`energyflowmonitor_spwr`**: Number
- **`energyflowmonitor_pre`**: Number
- **`energyflowmonitor_pri`**: Number
- **`energyflowmonitor_co2`**: Number
- **`energyflowmonitor_actual0`**: Number

### `Gate`
- **`gate_position`**: Number (0.0 - 1.0)
- **`gate_active`**: Number (-1, 0, 1)
- **`gate_preventopen`**: Boolean
- **`gate_preventclose`**: Boolean

### `Hourcounter`
- **`hourcounter_total`**: Number (Seconds)
- **`hourcounter_remaining`**: Number (Seconds)
- **`hourcounter_lastactivation`**: Number (Unix Timestamp)
- **`hourcounter_overdue`**: Boolean
- **`hourcounter_maintenanceinterval`**: Number (Seconds)
- **`hourcounter_active`**: Boolean

### `InfoOnlyAnalog`
- **`infoonlyanalog_value`**: Number

### `InfoOnlyDigital`
- **`infoonlydigital_value`**: Boolean
- **`infoonlydigital_text`**: String
- **`infoonlydigital_color`**: String

### `InfoOnlyText`
- **`infoonlytext_text`**: String

### `Intelligent Room Controller v2`
- **`intelligentroomcontrollerv2_activemode`**: Number (0-4)
- **`intelligentroomcontrollerv2_operatingmode`**: Number (0-5)
- **`intelligentroomcontrollerv2_tempactual`**: Number
- **`intelligentroomcontrollerv2_temptarget`**: Number
- **`intelligentroomcontrollerv2_comforttemperature`**: Number
- **`intelligentroomcontrollerv2_comforttemperaturecool`**: Number
- **`intelligentroomcontrollerv2_comforttolerance`**: Number
- **`intelligentroomcontrollerv2_openwindow`**: Boolean
- **`intelligentroomcontrollerv2_co2`**: Number
- **`intelligentroomcontrollerv2_humidityactual`**: Number

### `Intercom`
- **`intercom_bell`**: Boolean
- **`intercom_lastbellevents`**: String
- **`intercom_version`**: String

### `IntercomV2`
- **`intercomv2_bell`**: Boolean
- **`intercomv2_missedbellevents`**: String
- **`intercomv2_devicestate`**: Number

### `Jalousie`
- **`jalousie_up`**: Boolean
- **`jalousie_down`**: Boolean
- **`jalousie_position`**: Number (0.0 - 1.0)
- **`jalousie_shadeposition`**: Number (0.0 - 1.0)
- **`jalousie_safetyactive`**: Boolean
- **`jalousie_autoallowed`**: Boolean
- **`jalousie_autoactive`**: Boolean
- **`jalousie_locked`**: Boolean

### `LightControllerV2`
- **`lightcontrollerv2_activemoods`**: JSON Array (e.g., `[1, 4]`)
- **`lightcontrollerv2_moodlist`**: JSON Array
- **`lightcontrollerv2_favoritemoods`**: JSON Array
- **`lightcontrollerv2_additionalmoods`**: JSON Array
- **`lightcontrollerv2_presence`**: Boolean

### `Meter`
- **`meter_actual`**: Number
- **`meter_total`**: Number
- **`meter_totalneg`**: Number

### `NFC Code Touch`
- **`nfccodetouch_codedate`**: Number (Unix Timestamp)
- **`nfccodetouch_devicestate`**: Number (Bitmask)
- **`nfccodetouch_nfclearnresult`**: JSON Array

### `PresenceDetector`
- **`presencedetector_active`**: Boolean
- **`presencedetector_locked`**: Boolean
- **`presencedetector_activesince`**: Number (Unix Timestamp)

### `Pushbutton`
- **`pushbutton_active`**: Boolean

### `Radio`
- **`radio_activeoutput`**: Number (ID)

### `Remote`
- **`remote_mode`**: Number
- **`remote_active`**: Boolean

### `Sauna`
- **`sauna_active`**: Boolean
- **`sauna_power`**: Boolean
- **`sauna_tempactual`**: Number
- **`sauna_temptarget`**: Number
- **`sauna_fan`**: Boolean
- **`sauna_drying`**: Boolean
- **`sauna_error`**: Number
- **`sauna_saunaerror`**: Number

### `SmokeAlarm`
- **`smokealarm_nextlevel`**: Number
- **`smokealarm_level`**: Number
- **`smokealarm_acousticalarm`**: Boolean
- **`smokealarm_testalarm`**: Boolean
- **`smokealarm_alarmcause`**: Number (Bitmask)

### `SolarPumpController`
- **`solarpumpcontroller_buffertemp{n}`**: Number (Dynamic, n=0-4)
- **`solarpumpcontroller_bufferstate{n}`**: Number (Dynamic, n=0-4)
- **`solarpumpcontroller_collectortemp`**: Number

### `Switch`
- **`switch_active`**: Boolean (1=On, 0=Off)
- **`switch_jlocked`**: JSON Object (Lock info)
- **`switch_lockedon`**: Boolean

### `TextState`
- **`textstate_textandicon`**: String/JSON

### `TextInput`
- **`textinput_text`**: String

### `TimedSwitch`
- **`timedswitch_deactivationdelaytotal`**: Number (Seconds)
- **`timedswitch_deactivationdelay`**: Number (Seconds remaining)

### `Tracker`
- **`tracker_entries`**: String

### `ValueSelector`
- **`valueselector_value`**: Number
- **`valueselector_min`**: Number
- **`valueselector_max`**: Number
- **`valueselector_step`**: Number

### `Ventilation`
- **`ventilation_speed`**: Number (%)
- **`ventilation_mode`**: Number
- **`ventilation_temperatureindoor`**: Number
- **`ventilation_temperatureoutdoor`**: Number
- **`ventilation_temperaturetarget`**: Number
- **`ventilation_airqualityindoor`**: Number
- **`ventilation_humidityindoor`**: Number
- **`ventilation_presence`**: Boolean

### `Wallbox2`
- **`wallbox2_connected`**: Boolean
- **`wallbox2_enabled`**: Boolean
- **`wallbox2_active`**: Boolean
- **`wallbox2_limit`**: Number (kW)
- **`wallbox2_mode`**: Number
- **`wallbox2_session`**: JSON Object
- **`wallbox2_priceperkwh`**: Number

### `Window`
- **`window_position`**: Number (0.0 - 1.0)
- **`window_direction`**: Number (-1, 0, 1)
- **`window_lockedreason`**: String

### `WindowMonitor`
- **`windowmonitor_windowstates`**: String (List of states)
- **`windowmonitor_numopen`**: Number
- **`windowmonitor_numclosed`**: Number
- **`windowmonitor_numtilted`**: Number
- **`windowmonitor_numlocked`**: Number
- **`windowmonitor_numunlocked`**: Number


## `_info` Topics

These topics provide static metadata about the Miniserver, rooms, and controls. They are published as **retained** messages.

### Miniserver Info
**Topic:** `loxone/<serial>/_info`

Contains global configuration and Miniserver details.

```json
{
  "serialNr": "504F94A00000",
  "msName": "My Loxone Miniserver",
  "projectName": "My Smart Home",
  "localUrl": "192.168.1.200:80",
  "remoteUrl": "https://dns.loxonecloud.com/...",
  "tempUnit": 0,       // 0 = Celsius, 1 = Fahrenheit
  "currency": "€",
  "squareMeasure": "m²",
  "location": "Linz, Austria",
  "latitude": 48.3069,
  "longitude": 14.2858,
  "altitude": 260,
  "miniserverType": 2, // 0=Gen1, 1=Go Gen1, 2=Gen2, 3=Go Gen2, 4=Compact
  "catTitle": "Category",
  "roomTitle": "Room"
}
```

### Room Info
**Topic:** `loxone/<serial>/<room>/_info`

Metadata for a specific room. The `<room>` path segment is the normalized name (kebab-case) of the room.

```json
{
  "uuid": "10f3c6ba-0262-432d-8000959f23719000",
  "name": "Living Room",
  "image": "room_living_room.svg", // Icon name
  "defaultRating": 0
}
```

### Control Info
**Topic:** `loxone/<serial>/<room>/<control>/_info`

Metadata for a specific control.

```json
{
  "uuid": "10f3c6ba-0000-0000-0000000000000000",
  "name": "Ceiling Light",
  "type": "Switch",
  "cat": "Lighting",       // Name of the category
  "catUuid": "...",        // UUID of the category
  "room": "Living Room",   // Name of the room
  "roomUuid": "...",       // UUID of the room
  "isSecured": false,      // If visualization password is required
  "defaultRating": 0,
  "details": {             // Type-specific static details
    "allOff": "All Off",
    "outputs": {
      "1": "Radio 1"
    }
  }
}
```
