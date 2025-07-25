## stackit volume backup restore

Restores a backup

### Synopsis

Restores a backup by its ID.

```
stackit volume backup restore BACKUP_ID [flags]
```

### Examples

```
  Restore a backup with ID "xxx"
  $ stackit volume backup restore xxx
```

### Options

```
  -h, --help   Help for "stackit volume backup restore"
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none" "yaml"]
  -p, --project-id string      Project ID
      --region string          Target region for region-specific requests
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit volume backup](./stackit_volume_backup.md)	 - Provides functionality for volume backups

