# concat
concatenate a list of files to stdout with optional intermediate processing

## Usage

concat needs a yaml file to specify its actions. Files specified are processed in order. If a file has no recipe to match, the file will simply be read and output to stdout. 

Recipes work much like piped commands in that each command's stdout is connected to the next command's stdin. The first command recieves either the file name (ignoreFile: true) or file contents (ignoreFile: true) and the last command's stdout is routed to stdout.

Each recipe can specify an ```ignoreFile``` field. If ```ignoreFile``` is set to true concat does not try to open the filename as a file but instead passes it directly to the first command. This can for example be used wirh xargs and curl to download a remote file and insert it into stdout.

### Use Case

This program was originally written to be able to apply various local and remote patches in non-alphabetical order

## File Format
```
files:
  - <file 1>
  - <file 2>
  - <...>

recipes:
  - files:
    - <file 1>
    - <...>
    commands:
      - cmd: <command>
        args:
        - <arg 1>
        - <arg 2>
      - cmd: <command with no args>
  - files:
    - <url>
    commands:
      - cmd: xargs
        args:
          - curl
```
