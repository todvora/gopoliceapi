HOST: https://pcr-tdvorak.rhcloud.com/

# Stolen Vehicles Database API
Czech Police does not have any API for searching stolen vehicles. So i made one. 
If you need to verify registration number or VIN of vehicle, you can use my API. 

It's based on search form on oficial website [aplikace.policie.cz/patrani-vozidla](http://aplikace.policie.cz/patrani-vozidla/)  

More details and demo on [pcr-tdvorak.rhcloud.com](https://pcr-tdvorak.rhcloud.com/)

# Group Search in official Czech police database of stolen vehicles. 
There are used some typical abbreviations:

* __rpw__ - Registration Plates of the World - country od registration, see [en.wikipedia.org/wiki/List_of_international_vehicle_registration_codes](http://en.wikipedia.org/wiki/List_of_international_vehicle_registration_codes)
* __regno__ - Registration number - Licence plate number -[en.wikipedia.org/wiki/Vehicle_registration_plate](http://en.wikipedia.org/wiki/Vehicle_registration_plate)
* __vin__ - Vehicle Identification Number - [en.wikipedia.org/wiki/VIN](http://en.wikipedia.org/wiki/VIN)
* __stolendate__ - in format DD.MM.YYYY

## /search?q=8B67354
### GET
Search in official Czech police database of stolen vehicles. Parameters are:

* __vin__ - VIN of the vehicle (uses db index in the source database - much faster)
* __regno__ - registration number of the vehicle (uses db index in the source database - much faster)
* __q__ - In case you don't know if the query is regno or vin. Works for both of them, although much slower.


+ Response 200 (text/javascript)

        {
          "results": [
            {
              "class": "osobní vozidlo",
              "manufacturer": "VW",
              "type": "PASSAT VARIANT 3BGAWXX0SGSE WV",
              "color": "modrá tmavá",
              "regno": "8B67354",
              "rpw": "CZ",
              "vin": "WVWZZZ3BZ2E409524",
              "engine": "",
              "productionyear": "2002",
              "stolendate": "13.3.2013",
              "url": "http://aplikace.policie.cz/patrani-vozidla/Detail.aspx?id=368304"
            }
          ],
          "count": 1,
          "time": "2013-05-08T10:17:25.997Z"
        }
