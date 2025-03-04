package tileset

import (
	"encoding/json"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/qiekn/rpg-go/constants"
)

// every tileset must be able to give an image given an id
type Tileset interface {
	Img(id int) *ebiten.Image
}

////////////////////////////////////////////////////////////////////////
//                              Uniform                               //
////////////////////////////////////////////////////////////////////////

type UniformTilesetJson struct {
	Path string `json:"image"`
}

type UniformTileset struct {
	img *ebiten.Image
	gid int
}

func (u *UniformTileset) Img(id int) *ebiten.Image {
	id -= u.gid

	// get the position on the image where the tile id is
	srcX := id % 22
	srcY := id / 22

	// convert the src tile pos to pixel src position
	srcX *= constants.Tilesize
	srcY *= constants.Tilesize

	return u.img.SubImage(image.Rect(srcX, srcY, srcX+constants.Tilesize, srcY+constants.Tilesize)).(*ebiten.Image)
}

////////////////////////////////////////////////////////////////////////
//                              Dynamic                               //
////////////////////////////////////////////////////////////////////////

type TileJSON struct {
	Id     int    `json:"id"`
	Path   string `json:"image"`
	Width  int    `json:"imagewidth"`
	Height int    `json:"imageheight"`
}

type DynTilesetJSON struct {
	Tiles []*TileJSON `json:"tiles"`
}

type DynTileset struct {
	imgs []*ebiten.Image
	gid  int
}

func (d *DynTileset) Img(id int) *ebiten.Image {
	id -= d.gid
	return d.imgs[id]
}

////////////////////////////////////////////////////////////////////////
//                              Tileset                               //
////////////////////////////////////////////////////////////////////////

func NewTileset(path string, gid int) (Tileset, error) {
	// read file contents
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	{ // 1. return dynamic tileset
		if strings.Contains(path, "buildings") {
			var dynTilesetJSON DynTilesetJSON
			err = json.Unmarshal(contents, &dynTilesetJSON)
			if err != nil {
				return nil, err
			}

			// create dynamic tileset
			dynTileset := DynTileset{
				gid:  gid,
				imgs: []*ebiten.Image{},
			}

			// loop over tile data and load image for each
			for _, tileJSON := range dynTilesetJSON.Tiles {

				tileJSONPath := PathMagic(tileJSON.Path)

				img, _, err := ebitenutil.NewImageFromFile(tileJSONPath)
				if err != nil {
					return nil, err
				}

				dynTileset.imgs = append(dynTileset.imgs, img)
			}

			return &dynTileset, nil
		}
	}
	{ // 2. return uniform tileset
		var uniformTilesetJSON UniformTilesetJson
		err = json.Unmarshal(contents, &uniformTilesetJSON)
		if err != nil {
			return nil, err
		}

		uniformTileset := UniformTileset{}
		tileJSONPath := PathMagic(uniformTilesetJSON.Path)

		img, _, err := ebitenutil.NewImageFromFile(tileJSONPath)
		if err != nil {
			return nil, err
		}
		uniformTileset.img = img
		uniformTileset.gid = gid

		return &uniformTileset, nil
	}
}

func PathMagic(path_ string) (path string) {
	// convert tileset relative path to root relative path
	// original:  "..\/..\/images\/buildings\/building1.png"
	// output: "assets/images/buildings/building1.png"
	path = filepath.Clean(path_)
	path = strings.ReplaceAll(path, "\\", "/")
	path = strings.TrimPrefix(path, "../")
	path = strings.TrimPrefix(path, "../")
	path = filepath.Join("assets/", path)
	fmt.Println("magic path: ", path)
	return path
}
