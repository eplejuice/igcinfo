# Assignment 1: in-memory IGC track viewer

##### Martin Brådalen  |  martbraa@stud.ntnu.no  |  Studnr: 473145

## Assignment description
Develop an online service that will allow users to browse information about IGC files.

#### Assignment link: 
    http://prod3.imt.hig.no/teaching/imt2681-2018/wikis/assignment-1

#### Heroku link 
    https://vast-hamlet-56796.herokuapp.com/igcinfo/api
    
#### External dependencies
    https://github.com/marni/goigc
    https://github.com/gorilla/mux
 
### Quality
- [x] Golint
- [x] GoVet

## Test and expected results

### Tested using [Postman](https://www.getpostman.com/)

##### Bad Requests
    https://vast-hamlet-56796.herokuapp.com/        -> 404
    https://vast-hamlet-56796.herokuapp.com/igcinfo -> 404


#### GET /api
[Test](https://vast-hamlet-56796.herokuapp.com/igcinfo/api)

Returns metadata about the API

    {
      "uptime": uptime in the ISO8601 fomat
      "info": "Infomation about the API"
      "version": "1"
    }

#### POST /api/igc
[Test with postman](https://www.getpostman.com/)

Register a track, sent as a json in the body of the Url

        {
        "url": "http://skypolaris.org/wpcontent/uploads/IGS%20Files/Madrid%20to%20Jerez.igc"
        }

Returns the id assigned to the track

        {
          "id": "<id>"
        }

#### GET api/igc
[Test](https://vast-hamlet-56796.herokuapp.com/igcinfo/api/igc)

Returns an array with all IDs of all registered tracks.
Returns an empty array if there are no registered tracks.

    [<id1>, <id2>, ...]

#### GET api/igc/`<id>`
[Test with id = 1 ](https://vast-hamlet-56796.herokuapp.com/igcinfo/api/igc/1)

Returns info storen on a track, based on given `<id>`

Returns empty if no tracks are registered beforehand


        {
        "H_date": <date from File Header, H-record>,
        "pilot": <pilot>,
        "glider": <glider>,
        "glider_id": <glider_id>,
        "track_length": <calculated total track length>
        }

#### GET api/igc/`<id>`/`<field>`
[Test with id = 1 and field = "pilot"](https://vast-hamlet-56796.herokuapp.com/igcinfo/api/igc/1/pilot)

Returns a single value of a track, based on `<id>` and `<field>`

Returns empty if no tracks are registered beforehand

    <pilot> Miguel Angel Gordillo
