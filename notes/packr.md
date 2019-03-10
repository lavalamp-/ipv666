# Using Packr w/ IPv666 Assets
1. Create the file that should be packaged up with `ipv666`'s distribution
2. Use the `fs.ZLibCompress` function to create a `zlib` compressed copy of the file and place it in the `assets` directory
3. Add code to the `manager.go` file in the `data` module for retrieving the contents of the file and using them as necessary (see references to the `packedBox` variable for examples)
4. Run `packr2` from the root `ipv666` directory (this will create a new directory `packrd` with a file `packed-packr.go`)
5. Review the code in `packed-packr.go` and port over **only** the files you need to `internal/data/packed.go` (`packr` will attempt to package up files in the `data` module as they reference `packr` code)
6. Delete the `packrd` directory
7. Delete the `internal/data/data-packr.go` file
8. Test to make sure that the file is successfully loaded from `packed.go` in regular usage