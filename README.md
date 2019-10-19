# jsons2table
Serialise multiple same-format JSON files to a CSV and / or EXCEL file.

## Content

<!-- TOC depthFrom:2 depthTo:6 withLinks:1 updateOnSave:1 orderedList:0 -->

- [Content](#content)
- [Principles](#principles)
- [Installation](#installation)

<!-- /TOC -->

---
## Principles

Executing this command:

```sh
jsons2table /path/to/my/folder/with/json_files/my_folder_name
```

creates, within the given folder:

- an `Excel` (`.xlsx`) file with 1 line for each original `JSON` file
- also creates a `CSV` file with 1 line for each original `JSON` file
- if non-existent yet, a `.conf` file that is used to format the `Excel` file

[Top](#content)

---
## Installation

- make sure you have `Go` installed and working
- `go get -u github.com/ninjawule/jsons2table`

[Top](#content)
