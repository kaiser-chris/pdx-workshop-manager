package steam

// #cgo CFLAGS: -I${SRCDIR} -I${SRCDIR}/../sdk/public/steam
// #cgo LDFLAGS: -Wl,-rpath,:. -L${SRCDIR}/../sdk/redistributable_bin/linux64 -Wl,-Bdynamic -lsteam_api -static-libgcc -static-libstdc++
import "C"
