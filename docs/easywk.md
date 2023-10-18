# EasyWK Livetiming

> **Warning**
> For the use of SwimResults it is mandatory to deactivate the setting "Ready-Signal der Anlage schreibt neuen Lauf" in order to receive new heat data at the moment of the start and write the "start at" time correctly.

The import service is able to receive data from EasyWK Livetiming.
In a regular use case EasyWK issues HTTP POST Request against a php file (`livework.php`) which is then in charge of storing the received data and building HTML tables with the data to be displayed.
For reference the used `livework.php` file can be found in this directory. 

In SwimResults we are going to create an HTTP POST endpoint inside the import service that can be called by EasyWK in the same way as the php file.
This endpoint is then used to send the received data to the correct microservice.

To be able to create the endpoints it is mandatory to go through the `livework.php` file and analyse how the data is submitted by EasyWK.

## JSON Keywords

```text
'pwd'
'pwd'
'action'
'action'
'keepsum'
'firstlane'
'lanecount'
'vername'
'event'
'heat'
'maxheat'
'name'
'swr'.$i
'swr'.$i
'yob'.$i
'yob'.$i
'club'.$i
'club'.$i
'lane'
'meter'
'meter'
'time'
'meter'
'finished'
'finished'
'time'
'lane'
'content'
```

Inside the php file there are 28 locations where `$_REQUEST['keyword']` is used to read data from the given keyword.
The list above shows all used occurrences. Those with `$i` are using an iterator variable.

## Actions

Depending on the string defined inside the 'action' field of the received JSON different data will be submitted inside the request.
The following list includes all actions and the keywords that are used within this action's data.

If the action field is not set, the following text is returned by default: `ERROR: Passwort nicht korrekt oder keine Aktion definiert`

If the given action does not exist, the following text is returned by default: `ERROR: Unbekannte Aktion` 

### ping

simple ping

```json
{
  "pwd": "PASSWORD",
  "action": "ping"
}
```

`returns: 'OK'`

### clearsum

clears the summary file and list

```json
{
  "pwd": "PASSWORD",
  "action": "clearsum"
}
```

`returns: 'OK'`

### init

receives initial data for a meeting

```json
{
  "pwd": "PASSWORD",
  "action": "init",
  "keepsum": true,
  "firstlane": 1,
  "lanecount": 4,
  "vername": "Erzgebirgsschwimmcup"
}
```

`returns: 'OK'`

### newrace

receive athletes and heat data for the next heat; clear times

```json
{
  "pwd": "PASSWORD",
  "action": "newrace",
  "event": 13,
  "heat": 3,
  "maxheat": 15,
  "name": "50m Schmetterling männlich",
  "swr1": "Luca Heidenreich",
  "yob1": "2008",
  "club1": "ST Erzgebirge",
  "swr2": "Noah Joel Meusel",
  "yob2": "2008",
  "club2": "SV 1990 Zschopau",
  "swr3": "Alex Duckstein",
  "yob3": "2010",
  "club3": "SV 1919 Grimma",
  "swr4": "",
  "yob4": "",
  "club4": ""
}
```

`returns: 'OK'`

### ready

the same as `newrace`

```json
{
  "pwd": "PASSWORD",
  "action": "ready",
  "event": 13,
  "heat": 3,
  "maxheat": 15,
  "name": "50m Schmetterling männlich",
  "swr1": "Luca Heidenreich",
  "yob1": "2008",
  "club1": "ST Erzgebirge",
  "swr2": "Noah Joel Meusel",
  "yob2": "2008",
  "club2": "SV 1990 Zschopau",
  "swr3": "Alex Duckstein",
  "yob3": "2010",
  "club3": "SV 1919 Grimma",
  "swr4": "",
  "yob4": "",
  "club4": ""
}
```

`returns: 'OK'`

### time

receives a result time for a single lane. It is either "RT" for reaction time or a split or final result time. 

```json
{
  "pwd": "PASSWORD",
  "action": "time",
  "lane": 2,
  "meter": "RT",
  "time": 243240,
  "finished": "yes"
}
```

`returns: 'OK'`

### disq

sets the given lane to be disqualified, no reason submitted

```json
{
  "pwd": "PASSWORD",
  "action": "disq",
  "lane": 2
}
```

`returns: 'OK'`

### text

displays the submitted text in the livetiming

```json
{
  "pwd": "PASSWORD",
  "action": "text",
  "content": "Das ist ein Infotext!"
}
```

`returns: 'OK'`

If file includes scripts (JS or PHP):
`ERROR: Unerlaubte Zeichen im Text`

### raceresult

after a heat is finished (all athletes finished or manual finish) this action rewrites the livetiming table and changes the ordering to be ordered by rank instead of lane. 

```json
{
  "pwd": "PASSWORD",
  "action": "raceresult"
}
```

`returns: 'OK'`

## Password

With every request there is a password submitted by EasyWK. If this password is not correct, the following text is returned by default: `ERROR: Passwort nicht korrekt oder keine Aktion definiert`

## Import Service Behaviour

Depending on the given `action` several tasks have to be done by the import service.
The following table lists all actions and the corresponding import processes.

| Action     | Tasks                                                                    |
|------------|--------------------------------------------------------------------------|
| ping       | do nothing                                                               |
| clearsum   | do nothing                                                               |
| init       | do nothing                                                               |
| newrace    | set heat startet at time; remember current heat; remember current event; |
| ready      | set heat startet at time; remember current heat; remember current event; |
| time       | import result for given event, heat and lane                             |
| disq       | do nothing                                                               |
| text       | do nothing                                                               |
| raceresult | set heat finished at time                                                |


## Time Formats

Times are send by EasyWK as integer numbers. Reaction times have to be parsed as well as lap and finish times.

### Reaction Time

Reaction time is indicated in the "time" action by setting the `meters` field to `"RT"`.
The time consists of a number with 1 to 3 digits and is parsed by padding the number with zeros on the left until it reaches 3 digits length. After that a comma is added between the first and second digit.

```
345 -> "3,45"
23  -> "0,23"
3   -> "0,03"
```

### Result and Lap Times

Result and lap times are parsed in a similar way.
The number value is padded with zeros on the left until it reaches a length of 8 digits. This 8 digits are then inserted in the following pattern: `"00:00:00,00"`

```
12345678    -> "12:34:56,78"
2345678     -> "02:34:56,78"
345678      -> "00:34:56,78"
45678       -> "00:04:56,78"
5678        -> "00:00:56,78"
678         -> "00:00:06,78"
78          -> "00:00:00,78"
8           -> "00:00:00,08"
0           -> "00:00:00,00"
```
