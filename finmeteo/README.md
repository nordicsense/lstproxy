# finmeteo
A utility to mass-download Finnish meteo data from 
[opendata.fmi.fi](http://opendata.fmi.fi). The tool uses the 
`fmi::observations::weather::timevaluepair` query and permits downloading
temperature and other meteo data for given station IDs in a given time range.
The station IDs need to be provided as a CSV file where one can specify which
column to use for IDs and whether a certain number of header lines need to be 
skipped. 

The granularity of dates is 5 days, that is the app will split the given time
range into 5-day long intervals and will process each 5-days per station at a
time.
```
finmeteo [--from=string] [--to=string] [--params=string] [--col=int] [--skip=int] [--step=int] <stations> [path]

Description:
    Mass-downloader for Finnish meteo data

Arguments:
    stations       Stations CSV file
    path           Output path (default: working directory), optional

Options:
    -f, --from     From date as 2006-12-25 (default: -1y from today)
    -t, --to       Till date as 2006-12-25 (default: today)
    -p, --params   Parameters, coma separated w/o spaces (default: all)
    -c, --col      CSV column index (default: 2)
        --skip     CSV lines to skip (default: 1)
    -s, --step     Time step in minutes (default: 60)
```
