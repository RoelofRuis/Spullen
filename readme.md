## Spullen

A small application to allow me to index and search through all the things I own.

The aim of creating this index is becoming more aware of how many things I own and to try to reduce this set.


### Settings

The application is configured through passing of environment variables.

#### DBROOT
The application looks for existing databases in the same folder as the executable. To change this folder, set `DBROOT` to the appropriate root path.

#### PORT
The application is started on port `8080` by default. Change this by setting `PORT`.

#### MODE
Set `MODE` to `DEV` for development mode. This ensures that:
- The database files will then be output as plain text.
- The templates will be hot reloaded on each page refresh.