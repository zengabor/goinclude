# goinclude
Command-line build tool which includes the content of a file into a go template and outputs the merged result to stdout.

## Usage

    goinclude <file-path-of-template>

Example:

    $ goinclude templates/main.gohtml > templates/main.html

Where inside the `templates/main.gohtml` you reference other files via `{{ include ... }}`, e.g.,

    <style>{{ include "css/inline-styles.css" }}</style>

## Install

```bash
$ go get github.com/zengabor/goinclude
```

## Author

[Gabor Lenard](https://github.com/zengabor)
