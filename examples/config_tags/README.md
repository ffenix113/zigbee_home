## Tags that can be used in configuration file to make it more dynamic/configurable.


Tags included in this example:
* `!include` - tag can help split configuration into multiple files, allowing to re-use parts of configuration in multiple files files.
* `!env` - fetch value from environment and replace the tag with it. If env variable is not present - configuration may become invalid and might break the build process.


Env variable `ZBHOME_BOARD_NAME` will be used to set the board name, and file `board.yaml` contains board definition that will be included inside the main config file.