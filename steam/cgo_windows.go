package steam

// #cgo CFLAGS: -I${SRCDIR} -I${SRCDIR}/../sdk/public/steam
// #cgo LDFLAGS: -Wl,-rpath,:. -L${SRCDIR}/../sdk/redistributable_bin/win64 -Wl,-Bdynamic -lsteam_api64 -static-libgcc -static-libstdc++
import "C"
