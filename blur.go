package main

/*#cgo LDFLAGS: -lxcb
#cgo CFLAGS: -w
#include <xcb/xcb.h>
#include <xcb/xproto.h>
#include <string.h>
#include <stdlib.h>

xcb_atom_t GetAtom(xcb_connection_t *conn, char *name){
	xcb_atom_t atom;
	xcb_intern_atom_cookie_t cookie;

	cookie = xcb_intern_atom(conn, 0, strlen(name), name);

	xcb_intern_atom_reply_t *reply = xcb_intern_atom_reply(conn, cookie, NULL);
	if(reply) {
		atom = reply->atom;
		free(reply);
	}
	return atom;
}

int GetWindows(xcb_connection_t *conn) {
	xcb_atom_t atom = GetAtom(conn, "_NET_CLIENT_LIST");
	xcb_get_property_cookie_t prop_cookie;
	prop_cookie = xcb_get_property(conn, 0, xcb_setup_roots_iterator(xcb_get_setup(conn)).data->root, atom, 0, 0, (1 << 32)-1);
	xcb_get_property_reply_t *prop_reply;
	prop_reply = xcb_get_property_reply(conn, prop_cookie, NULL);

	return prop_reply->value_len;
}
*/
import "C"

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func MakeGaussian(size int, max float64) []float64 {
	var ret []float64

	for x := float64(-1); x <= 1; x += float64(1 / (float64(size) - 1) * 2) {
		Calc := math.Pow(math.E, -(math.Pi * math.Pow(x, 2)))
		ret = append(ret, Calc)
	}

	ratio := max / 1
	for i, _ := range ret {
		ret[i] = math.Round(ratio * ret[i])
		if ret[i] == 0 {
			ret[i] = 1
		}
	}

	return ret
}

func PutMask(imgIn image.Image, mask []float64) image.Image {
	imgOut := image.NewRGBA(imgIn.Bounds())
	bounds := imgIn.Bounds()
	var img = make([][3]uint32, (bounds.Min.Y-bounds.Max.Y)*(bounds.Min.X-bounds.Max.X))
	for y := 0; y < bounds.Max.Y-bounds.Min.Y; y++ {
		for x := 0; x < bounds.Max.X-bounds.Min.X; x++ {
			r, g, b, _ := imgIn.At(x, y).RGBA()
			colors := [3]uint32{r / 257, g / 257, b / 257}
			img[x+(y*(bounds.Max.X-bounds.Min.X))] = colors
		}
	}

	for y := 0; y < bounds.Max.Y-bounds.Min.Y; y++ {
		for x := 0; x < bounds.Max.X-bounds.Min.X; x++ {
			newR, newB, newG := 0, 0, 0
			sum := 0
			for i := -(len(mask) / 2); i < len(mask)/2; i++ {
				if x+i >= bounds.Max.X || x+i < 0 {
					continue
				}
				sum += int(mask[i+len(mask)/2])
				index := x + i + (y * (bounds.Max.X - bounds.Min.X))
				r, g, b := img[index][0], img[index][1], img[index][2]
				newR += int(mask[i+len(mask)/2]) * int(r)
				newG += int(mask[i+len(mask)/2]) * int(g)
				newB += int(mask[i+len(mask)/2]) * int(b)
			}
			newR /= sum
			newG /= sum
			newB /= sum
			img[x+(y*bounds.Max.X-bounds.Min.X)] = [3]uint32{uint32(newR), uint32(newG), uint32(newB)}
		}
	}
	for y := 0; y < bounds.Max.Y-bounds.Min.Y; y++ {
		for x := 0; x < bounds.Max.X-bounds.Min.X; x++ {
			newR, newB, newG := 0, 0, 0
			sum := 0
			for i := -(len(mask) / 2); i < len(mask)/2; i++ {
				if y+i >= bounds.Max.Y || y+i < 0 {
					continue
				}
				sum += int(mask[i+len(mask)/2])
				index := x + ((y + i) * (bounds.Max.X - bounds.Min.X))
				r, g, b := img[index][0], img[index][1], img[index][2]
				newR += int(mask[i+len(mask)/2]) * int(r)
				newG += int(mask[i+len(mask)/2]) * int(g)
				newB += int(mask[i+len(mask)/2]) * int(b)
			}
			newR /= sum
			newG /= sum
			newB /= sum
			img[x+(y*bounds.Max.X-bounds.Min.X)] = [3]uint32{uint32(newR), uint32(newG), uint32(newB)}
		}
	}

	for y := 0; y < bounds.Max.Y-bounds.Min.Y; y++ {
		for x := 0; x < bounds.Max.X-bounds.Min.X; x++ {
			index := x + (y * (bounds.Max.X - bounds.Min.X))
			imgOut.Set(x, y, color.RGBA{uint8(img[index][0]), uint8(img[index][1]), uint8(img[index][2]), 0xff})
		}
	}

	return imgOut
}

func SetWallpaper(command string, filename string) error {
	s := strings.Fields(command + " " + filename)
	fmt.Println(s)
	if err := exec.Command(s[0], s[1:]...).Run(); err != nil {
		return err
	}
	return nil
}

func main() {
	increment := 5
	size := 5
	if true {
		fmt.Println("Creating Images")
		imagePath := os.Args[1]

		data, err := os.Open(imagePath)
		if err != nil {
			panic(err)
		}

		imgIn, _, err := image.Decode(data)
		if err != nil {
			panic(err)
		}

		for i := 0; i < size; i++ {
			start := time.Now()
			size := i * increment
			if size%2 == 0 || size == 0 {
				size++
			}
			mask := MakeGaussian(size, float64(i*increment))
			fmt.Println(len(mask))

			if i == 0 {
				mask = []float64{0, 1}
			}

			filename := strconv.Itoa(i) + "b.png"
			f, _ := os.Create(filename)
			png.Encode(f, PutMask(imgIn, mask))
			fmt.Println(i, " image processed :", time.Now().Sub(start), " Using mask :", mask)
		}
	}

	XConn := C.xcb_connect(nil, nil)
	defer C.xcb_disconnect(XConn)
	BlurOrder := make(chan bool)

	// Blurer
	go func() {
		blurStatus := false
		for {
			blur := <-BlurOrder
			if blur == blurStatus {
				continue
			} else if blur {
				for i := 1; i < size; i++ {
					SetWallpaper("feh --bg-fill", strconv.Itoa(i)+"b.png")
					//time.Sleep(time.Millisecond)
				}
				blurStatus = true
			} else if !blur {
				for i := size - 2; i >= 0; i-- {
					SetWallpaper("feh --bg-fill", strconv.Itoa(i)+"b.png")
					//time.Sleep(time.Millisecond)
				}
				blurStatus = false
			}
		}
	}()
	for {
		if C.GetWindows(XConn) > 1 {
			BlurOrder <- true
		} else {
			BlurOrder <- false
		}
	}
}
