# jsons2table
Serialise multiple same-format JSON files to a CSV and / or EXCEL file.

## Content

<!-- TOC depthFrom:2 depthTo:6 withLinks:1 updateOnSave:1 orderedList:0 -->

- [Content](#content)
- [Principles](#principles)
- [Installation](#installation)
- [License](#license)

<!-- /TOC -->

---
## Principles

Executing this command:

```sh
jsons2table /path/to/my/folder/with/json_files/my_folder_name
```

creates, within the given folder:

- an `Excel` (`.xlsx`) file with 1 line for each original `JSON` file: `my_folder_name.xlsx`
- also creates a `CSV` file with 1 line for each original `JSON` file: `my_folder_name.csv`
- if non-existent yet, a `.conf` file that is used to format the `Excel` file: `my_folder_name.conf`

**NB**: the previous version of the `Excel` and `CSV` files are erased during the process, so be careful.

[Top](#content)

---
## Installation

- make sure you have `Go` installed and working
- `go get -u github.com/ninjawule/jsons2table`

[Top](#content)

---
## License

This program is under the terms of the [MIT License](LICENSE).

It uses code from [tealeg/xlsx](https://github.com/tealeg/xlsx), which is [BSD 3-Clause](https://github.com/tealeg/xlsx/blob/master/LICENSE)-licensed.
