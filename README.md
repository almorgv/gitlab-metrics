# gitlab-metrics

Collect projects events from Gitlab and store it in PostgreSQL

## Configuration via environment

- `DB_HOST` - Database host
- `DB_PORT` - Database port
- `DB_USER` - Database user
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name
- `GITLAB_URL` - URL of gitlab instance
- `GITLAB_TOKEN` - Gitlab access token (requires `read_api` privileges)
- `UPDATE_INTERVAL` - Interval to fetch new events in minutes
- `LOG_LEVEL` - Logging level (`trace|debug|info|warn|error` - default `info`)
- `LOG_MODE` - Logging mode (`pretty`|`json` - default `pretty`)
