### Flashing
CLI can flash already built application with a couple of methods:
- nrfutil
- mcuboot
- west
- adafruit (if available in current PATH)

Device already has to be in a mode that will allow flashing(DFU, for example).

Example:
```sh
go run ./cmd/zigbee --config ./zigbee_test.yml firmware --workdir ./firmware flash
```

Also see [full example](index.md#full-example) for flashing instructions.