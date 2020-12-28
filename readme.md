## Spullen

A small application to allow me to index and search through all the things I own.

The aim of creating this index is becoming more aware of how many things I own and to try to reduce this set.


### Settings

#### DBROOT

The application looks for existing databases in the same folder as the executable. To change this folder, set the env variable `DBROOT` to the appropriate root path.

#### PORT
The application is started by default on port `8080`. Change this by setting the `PORT` env variable.

#### MODE
Set `MODE` to `DEV` for development mode. This ensures that:
- The database files will then be output as plain text.
- The templates will be hot reloaded on each page refresh.