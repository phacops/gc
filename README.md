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

## Create a configuration file (optional)

`gc` will try to locate a configuration file at
```
${XDG_CONFIG_HOME}/gc/config
${HOME}/.config/gcrc
${HOME/.gcrc
```

It should be a valid json file. Valid keys are `gc_username`, `gc_password`, and
`watch_dir`. The first two variables are Garmin Connect credentials, the last
one is where the Garmin device is mounted. If some of these variables are not
set, the user will be prompted to type them.

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

Without config file:
```
$ gc sync activities
Garmin Connect Username: user@example.com
Garmin Connect Password: 
Watch Mount Directory: /mnt
syncing ACTIVITY1.FIT... success
```

If you saved your settings in a config file, then no prompt is issued:
```
$ gc sync activities
syncing ACTIVITY3.FIT... success
```

Username and watch mount directory can be overriden with the `--username, -u`
and `--dir, -d` options respectively.

You can get more info using the built-in help:
```
$ gc --help
$ gc sync --help
```

## License

`gc` is released under the MIT license. Please see LICENSE file.
