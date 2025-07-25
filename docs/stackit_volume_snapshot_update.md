## stackit volume snapshot update

Updates a snapshot

### Synopsis

Updates a snapshot by its ID.

```
stackit volume snapshot update SNAPSHOT_ID [flags]
```

### Examples

```
  Update a snapshot name with ID "xxx"
  $ stackit volume snapshot update xxx --name my-new-name

  Update a snapshot labels with ID "xxx"
  $ stackit volume snapshot update xxx --labels key1=value1,key2=value2
```

### Options

```
  -h, --help                    Help for "stackit volume snapshot update"
      --labels stringToString   Key-value string pairs as labels (default [])
      --name string             Name of the snapshot
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

* [stackit volume snapshot](./stackit_volume_snapshot.md)	 - Provides functionality for snapshots

