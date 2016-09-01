# gc

gc is a command-line tool to sync activities, workouts, and so on to the Garmin
Connect service. It is *NOT* an official tool. It mainly targets the Linux as
Garmin doesn't provide the Connect Launcher for this platform.

## Installation

```
git clone https://github.com/phacops/gc.git
cd gc
go install
```

## Create a configuration file

`gc` will try to locate a configuration file at
```
${XDG_CONFIG_HOME}/gc/config
${HOME}/.config/gcrc
${HOME/.gcrc
```
should be a valid json file with `gc_username`, `gc_password`, and `watch_dir`
set. The first two variables are Garmin Connect credentials, the last one is
where the Garmin device is mounted.

Example:
```
{
  "gc_username":"user@example.com",
  "gc_password":"aVeryComplexPassword",
  "watch_dir":"/mnt"
}
```

## Usage

`gc` supports syncing activities, workouts, and download updated EPO file. The
Garmin device must first be mounted.

```
$ gc sync activities
syncing ACTIVITY3.FIT... success
```

You can get more info using the built-in help:
```
$ gc --help
$ gc sync --help
```

## License

`gc` is released under the MIT license. Please see LICENSE file.
