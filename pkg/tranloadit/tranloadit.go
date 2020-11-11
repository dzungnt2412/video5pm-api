package transloadit

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"video5pm-api/core/constants"

	transloadit "gopkg.in/transloadit/go-sdk.v1"
)

func CreatLayoutPreview(audioUrl string, length string, subPath string) string {
	options := transloadit.DefaultConfig
	options.AuthKey = constants.TRANLOADIT_KEY
	options.AuthSecret = constants.TRANLOADIT_SECRET
	client := transloadit.NewClient(options)

	// Initialize new Assembly
	assembly := transloadit.NewAssembly()

	// Set Encoding Instructions
	assembly.AddStep(":original", map[string]interface{}{
		"robot": "/upload/handle",
	})

	assembly.AddStep("resized", map[string]interface{}{
		"use":               ":original",
		"robot":             "/image/resize",
		"result":            true,
		"height":            768,
		"imagemagick_stack": "v2.0.7",
		"resize_strategy":   "fit",
		"width":             1024,
		"background":        "#000000",
	})

	assembly.AddStep("import_audio", map[string]interface{}{
		"robot": "/http/import",
		"url":   audioUrl,
	})

	mergedJson := `{
		"use": {"steps":[{"name":"import_audio", "as":"audio"},{"name":"resized", "as":"image"}],"bundle_steps":true},
		"robot": "/video/merge",
		"result": true,
		"duration": ` + length + `,
		"ffmpeg_stack": "v3.3.3",
		"framerate": "1",
		"preset": "ipad-high",
		"resize_strategy": "fit"
	  }`
	assembly.AddStep("merged", json2Map(mergedJson))

	assembly.AddStep("video_filtered", map[string]interface{}{
		"use":    "merged",
		"robot":  "/file/filter",
		"result": true,
		"accepts": []interface{}{
			[]interface{}{"${file.mime}", "regex", "video"},
		},
		"error_on_decline": false,
	})

	assembly.AddStep("subtitle_filtered", map[string]interface{}{
		"use":    ":original",
		"robot":  "/file/filter",
		"result": true,
		"accepts": []interface{}{
			[]interface{}{"${file.mime}", "regex", "application/x-subrip"},
		},
		"error_on_decline": false,
	})

	subtitledJson := `{
		"use": {"bundle_steps":true,"steps":[{"name":"video_filtered", "as":"video"},{"name":"subtitle_filtered", "as":"subtitles"}]},
		"robot": "/video/subtitle",
		"ffmpeg_stack": "v3.3.3",
		"preset": "ipad-high"
	  }`

	assembly.AddStep("subtitled", json2Map(subtitledJson))

	// mergedAudioAndVideoJson := `{
	// 	"use": {"steps":[{"name":"concatenated", "as":"audio"},{"name":"merged", "as":"video"}],"bundle_steps":true},
	// 	"robot": "/video/merge",
	// 	"result": true,
	// 	"ffmpeg_stack": "v3.3.3",
	// 	"preset": "ipad-high"
	//   }`

	// assembly.AddStep("mergedAudioAndVideo", json2Map(mergedAudioAndVideoJson))

	// exportedJson := `{"use": ["imported_chameleon", "imported_prinsengracht", "imported_snowflake", "resized", "merged", ":original"],"robot": "/dropbox/store","credentials": "Dropbox_Credentials"}`
	// assembly.AddStep("exported", json2Map(exportedJson))

	// Add files to upload
	assembly.AddFile("myfile_1", "../../public/layout.png")
	assembly.AddFile("myfile_2", subPath)

	// Start the Assembly
	info, err := client.StartAssembly(context.Background(), assembly)
	if err != nil {
		panic(err)
	}

	// All files have now been uploaded and the Assembly has started but no
	// results are available yet since the conversion has not finished.
	// WaitForAssembly provides functionality for polling until the Assembly
	// has ended.
	info, err = client.WaitForAssembly(context.Background(), info)
	if err != nil {
		panic(err)
	}

	fmt.Printf("You can check some results at: \n")
	fmt.Println("  - %s\n", info)
	// fmt.Println("  - %s\n", info.Results["merged"][0].SSLURL)
	// fmt.Printf("  - %s\n", info.Results["exported"][0].SSLURL)

	return info.Results["subtitled"][0].SSLURL

}

func CreateAudio(listAudio []string) string {
	options := transloadit.DefaultConfig
	options.AuthKey = constants.TRANLOADIT_KEY
	options.AuthSecret = constants.TRANLOADIT_SECRET
	client := transloadit.NewClient(options)

	// Initialize new Assembly
	assembly := transloadit.NewAssembly()

	// Set Encoding Instructions
	assembly.AddStep(":original", map[string]interface{}{
		"robot": "/upload/handle",
	})

	fileName := ""
	path := ""

	for i, v := range listAudio {

		fileName = "myfile_" + strconv.Itoa(i+1)
		path = "../../public/audio/" + v

		assembly.AddFile(fileName, path)
		fmt.Println(i)
		time.Sleep(1 * time.Second)
	}

	concatenatedJson := `{
		"use": {"steps":[{"name":":original", "as":"audio"}]},
		"robot": "/audio/concat",
		"result": true,
		"ffmpeg_stack": "v3.3.3"
	  }`

	fmt.Println(concatenatedJson)

	assembly.AddStep("concatenated", json2Map(concatenatedJson))

	// Start the Assembly
	info, err := client.StartAssembly(context.Background(), assembly)
	if err != nil {
		panic(err)
	}

	// All files have now been uploaded and the Assembly has started but no
	// results are available yet since the conversion has not finished.
	// WaitForAssembly provides functionality for polling until the Assembly
	// has ended.
	info, err = client.WaitForAssembly(context.Background(), info)
	if err != nil {
		panic(err)
	}

	fmt.Printf("You can check some results at: \n")
	fmt.Println("  - %s\n", info)
	fmt.Println("  - %s\n", info.Results["concatenated"][0].SSLURL)
	// fmt.Println("  - %s\n", info.Results["merged"][0].SSLURL)
	// fmt.Printf("  - %s\n", info.Results["exported"][0].SSLURL)

	return info.Results["concatenated"][0].SSLURL
}

func Concatenate_video(video_sentence string, video_sentence_length time.Duration, video_preview string, video_preview_length time.Duration, length time.Duration, audio string, sub string) (string, string, int64) {
	options := transloadit.DefaultConfig
	options.AuthKey = constants.TRANLOADIT_KEY
	options.AuthSecret = constants.TRANLOADIT_SECRET
	client := transloadit.NewClient(options)
	// Initialize new Assembly
	assembly := transloadit.NewAssembly()

	// Set Encoding Instructions
	assembly.AddStep(":original", map[string]interface{}{
		"robot": "/upload/handle",
	})

	assembly.AddStep("resized", map[string]interface{}{
		"use":               ":original",
		"robot":             "/image/resize",
		"result":            true,
		"height":            768,
		"imagemagick_stack": "v2.0.7",
		"resize_strategy":   "fit",
		"width":             1024,
		"background":        "#000000",
	})

	lengthLayout := length - video_preview_length - video_sentence_length

	fmt.Println(int64(video_sentence_length.Milliseconds()))
	fmt.Println(int64(video_preview_length.Milliseconds()))
	fmt.Println(int64(length.Milliseconds()))

	if int64(lengthLayout.Milliseconds()) > 500 {
		mergedJson := `{
			"use": {"steps":[{"name":"resized", "as":"image"}],"bundle_steps":true},
			"robot": "/video/merge",
			"result": true,
			"duration": ` + strconv.FormatInt(int64(lengthLayout.Seconds()), 10) + `,
			"ffmpeg_stack": "v3.3.3",
			"framerate": "1",
			"preset": "ipad-high",
			"resize_strategy": "fit"
		  }`
		assembly.AddStep("merged", json2Map2(mergedJson, 1))

		assembly.AddStep("merged_resized", map[string]interface{}{
			"use":             "merged",
			"robot":           "/video/encode",
			"result":          true,
			"background":      "#000000",
			"ffmpeg_stack":    "v3.3.3",
			"height":          270,
			"preset":          "ipad-high",
			"resize_strategy": "pad",
			"width":           480,
		})
	}

	if video_preview != "" {
		assembly.AddStep("preroll_imported", map[string]interface{}{
			"robot":  "/http/import",
			"result": true,
			"url":    video_preview,
		})

		assembly.AddStep("preroll_resized", map[string]interface{}{
			"use":             "preroll_imported",
			"robot":           "/video/encode",
			"result":          true,
			"background":      "#000000",
			"ffmpeg_stack":    "v3.3.3",
			"height":          270,
			"preset":          "ipad-high",
			"resize_strategy": "pad",
			"width":           480,
		})
	}

	assembly.AddStep("emptysound", map[string]interface{}{
		"robot":  "/http/import",
		"result": true,
		"url":    "http://demos.transloadit.com/inputs/empty-sound.mp3",
	})

	removeAudioJson := `{
		"use": {"steps":[{"name":"emptysound", "as":"audio"},{"name":":original", "as":"video"}]},
		"robot": "/video/merge",
		"result": true,
		"duration": ` + strconv.FormatInt(int64(video_sentence_length.Seconds()), 10) + `,
		"ffmpeg_stack": "v3.3.3",
		"preset": "ipad-high"
	  }`

	assembly.AddStep("remove_udio", json2Map2(removeAudioJson, 2))

	assembly.AddStep("original_resized", map[string]interface{}{
		"use":             "remove_udio",
		"robot":           "/video/encode",
		"result":          true,
		"background":      "#000000",
		"ffmpeg_stack":    "v3.3.3",
		"height":          270,
		"preset":          "ipad-high",
		"resize_strategy": "pad",
		"width":           480,
	})

	if video_preview != "" {
		concatPreviewJson := `{"use": {"steps":[{"name":"preroll_resized", "as":"video_1"},{"name": "original_resized", "as":"video_2"}]},"robot": "/video/concat","result": true,"ffmpeg_stack": "v3.3.3","preset": "ipad-high"}`
		assembly.AddStep("concatPreview", json2Map2(concatPreviewJson, 3))

		concatenatedJson := `{"use": {"steps":[{"name":"concatPreview", "as":"video_1"},{"name": "merged_resized", "as":"video_2"}]},"robot": "/video/concat","result": true,"ffmpeg_stack": "v3.3.3","preset": "ipad-high"}`
		assembly.AddStep("concatenated", json2Map2(concatenatedJson, 4))
	} else {
		concatFinalJson := `{"use": {"steps":[{"name":"original_resized", "as":"video_1"},{"name": "merged_resized", "as":"video_2"}]},"robot": "/video/concat","result": true,"ffmpeg_stack": "v3.3.3","preset": "ipad-high"}`
		assembly.AddStep("concatenated", json2Map2(concatFinalJson, 5))
	}

	assembly.AddStep("import_audio", map[string]interface{}{
		"robot": "/http/import",
		"url":   audio,
	})

	mergedAudioAndVideoJson := `{"use": {"steps":[{"name":"import_audio", "as":"audio"},{"name":"concatenated", "as":"video"}]},"robot": "/video/merge","result": true,	"duration": ` + strconv.FormatInt(int64(length.Seconds()), 10) + `,"ffmpeg_stack": "v3.3.3","preset": "ipad-high"}`
	assembly.AddStep("mergedAudioAndVideo", json2Map2(mergedAudioAndVideoJson, 6))

	assembly.AddStep("video_filtered", map[string]interface{}{
		"use":    "mergedAudioAndVideo",
		"robot":  "/file/filter",
		"result": true,
		"accepts": []interface{}{
			[]interface{}{"${file.mime}", "regex", "video"},
		},
		"error_on_decline": false,
	})

	assembly.AddStep("subtitle_filtered", map[string]interface{}{
		"use":    ":original",
		"robot":  "/file/filter",
		"result": true,
		"accepts": []interface{}{
			[]interface{}{"${file.mime}", "regex", "application/x-subrip"},
		},
		"error_on_decline": false,
	})

	subtitledJson := `{"use": {"bundle_steps":true,"steps":[{"name":"video_filtered", "as":"video"},{"name":"subtitle_filtered", "as":"subtitles"}]},"robot": "/video/subtitle","ffmpeg_stack": "v3.3.3","preset": "ipad-high"}`

	assembly.AddStep("subtitled", json2Map2(subtitledJson, 7))

	var lengthPreview int64 = 0

	if video_preview != "" {
		lengthPreview = int64(video_sentence_length.Milliseconds()) + int64(video_preview_length.Milliseconds())
	} else {
		lengthPreview = int64(video_sentence_length.Milliseconds())
	}

	// Add files to upload
	assembly.AddFile("myfile_1", video_sentence)
	assembly.AddFile("myfile_2", "../../public/layout.png")
	assembly.AddFile("myfile_3", sub)

	// Start the Assembly
	info, err := client.StartAssembly(context.Background(), assembly)
	if err != nil {
		panic(err)
	}

	// All files have now been uploaded and the Assembly has started but no
	// results are available yet since the conversion has not finished.
	// WaitForAssembly provides functionality for polling until the Assembly
	// has ended.
	info, err = client.WaitForAssembly(context.Background(), info)
	if err != nil {
		panic(err)
	}

	fmt.Printf("You can check some results at: \n")
	// fmt.Printf("  - %s\n", info.Results["preroll_imported"][0].SSLURL)
	// fmt.Printf("  - %s\n", info.Results["original_resized"][0].SSLURL)

	if video_preview != "" {
		return info.Results["subtitled"][0].SSLURL, info.Results["concatPreview"][0].SSLURL, lengthPreview
	}

	return info.Results["subtitled"][0].SSLURL, info.Results["original_resized"][0].SSLURL, lengthPreview

}

func json2Map(jSon string) map[string]interface{} {
	byt := []byte(jSon)

	var dat map[string]interface{}

	if err := json.Unmarshal(byt, &dat); err != nil {
		panic(err)
	}

	return dat
}

func json2Map2(jSon string, i int) map[string]interface{} {
	byt := []byte(jSon)

	var dat map[string]interface{}

	if err := json.Unmarshal(byt, &dat); err != nil {
		fmt.Println(i)
		panic(err)
	}

	return dat
}
