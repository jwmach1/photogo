# Photo Go
A better approach to extracting your photos from Google to your personal cloud. I'm moving my photos out of Google to a Synology NAS.

* create a directory structure to organize the media
* each file has the correct timestamp of the media

## Why do this
Google Takeout is provided for getting all of your own media -- photos, videos -- out of Google. 
The timestamp of those files will not match the metadata in the media.Your Google Takout will have the timestamp of when you extract the zip. Also, there is no directory structure -- you'll have one directory with all the files.

Two choices:
* Extract the Takout zip files, then open the image files, and change the timestamp of the files
* use an API to pull the images and write the with the correct timestamp

### Use Google Takeout?
Do you have a lot of disk space to pull down gigs and terrabytes of zip files and extract those zip files? Well, Your setting up a NAS, so sure.  Extract those image files into your network storage and you have one big directory with filenames dated of when you extracted them.  Not likely what you want.  Synology Photo displays the photos based on their file timestamp, not the image metadata. Now your photo timeline is all dorked up.

### Use Google REST API?
We can spend a weekend writing a little program to use the REST Api to meet the goals. Synology NAS, and probably your NAS, support an SMB or other mount to your "home" directory within the NAS. For me this was `/Volumes/home/Photos/`.This little utility can write the files directly onto the NAS as the media is pulled via the REST api.

I took the initial step to setup the NAS photo library application on my phone, and turned on backup.  That gives me directory structure to model. But my phone does not have all my photo's going back decades.

## Setup
We use the Google Photos Library.

* start with https://developers.google.com/photos/library/guides/get-started
  * BUT, when you create the OAUTH certificate, use desktop instead of web
* download the oauth certificate json, and save it to `credentials.json` in the root folder here (note that the filename is expected and included in the `.gitignore`).
* Mount your NAS directory. for me this was in `/Volumes/home`. Use whatever you want.

## Running
There are two optional arguments:
* output -- the base/root directory where the media will be saved
* worker-count -- how many "workers" will be used to call the REST api

Pass your own output directory based on your NAS mounted path
> go run main.go -output "/Volumes/home/Photos/..."
 

 ## Verification
 Did it work!? Can you delete your Google Photos with the confidence that every memory is in your private cloud?
 ```sh
  find /Volumes/home/Photos -size 0
  ```

Did the final count printed at the end match _about_ that shown in [your google dashboard](https://myaccount.google.com/dashboard)?

Open your new photo library (Synology Photos?) and look for pictures at the top/newest that shouldn't be there.  They did not get the file creation time modification correctly.  My solution was to just delete them, and run the whole thing again. It only takes a few minutes to process 4k media files.

Finally, have the operating system help you verify the media came over correctly by looking at the `file` output for each:
```sh
find /Volumes/home/Photos/MobileBackup/Pixel\ Tré/DCIM/Camera -type f > files
file -f files >file_type_output
```

I then used vim and removed all the known types -- jpeg, and movies -- from the `file_type_output`.  The only one that didn't below was a random MacOS `.DSStore` file where I must have opened finder and previewed the images.

## Warning

Ultimately it is up to you if you have the confidence in your new backup complete.  Can you delete your photos and videos from Google Photos with the knowledge that your media is safe? Only you can decide.  However, this backup does not include Albums. This did not recreate your shared Albums.