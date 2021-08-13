# gcs-proxy

gcs-proxy is a tool that acts as a server for your GCS bucket.

This project is forked from [NYTimes/gcs-helper](https://github.com/nytimes/gcs-helper). Since the original repository is archived, we have forked it to maintain and further improve it.

It was designed to be used with [Kaltura's nginx-vod-module](https://github.com/kaltura/nginx-vod-module), but it can also work as a standalone proxy.

Specific to nginx-vod-module, gcs-helper provides support for the mapped mode
(when using the proper environment variables - ``GCS_HELPER_PROXY_PREFIX``,
``GCS_HELPER_MAP_PREFIX`` and ``GCS_HELPER_MAP_REGEX_FILTER``).

## Configuration

The following environment variables control the behavior of gcs-helper:

| Variable                         | Default value | Required | Description                                                                                                  |
| -------------------------------- | ------------- | -------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| GCS_HELPER_LISTEN                | :8080         | No       | Address to bind the server                                                                                                                                               |
| GCS_HELPER_BUCKET_NAME           |               | Yes      | Name of the bucket                                                                                                                                                       |
| GCS_HELPER_ROOT_REDIRECT         |               | No       | Redirect all requests from `/` to this path
| GCS_HELPER_LOG_LEVEL             | debug         | No       | Logging level                                                                                                                                                           |
| GCS_HELPER_PROXY_PREFIX          |               | No       | Prefix to use for the proxy binding. Required if running in map and proxy modes (example value: ``/proxy/``)                                                        |
| GCS_HELPER_PROXY_TIMEOUT         | 10s           | No       | Defines the maximum time in serving the proxy requests, this is a hard timeout and includes retries                                                                    |
| GCS_HELPER_MAP_PREFIX            |               | No       | Prefix to use for the map binding. Required if running in map and proxy modes (example value: ``/map/``)                                                                |
| GCS_HELPER_MAP_REGEX_FILTER      |               | No       | A regular expression that is used to deliver only those files that match the specified naming convention (example value: \d{3,4}p(\.mp4|[a-z0-9_-]{37}\.(vtt|srt))$) |

The are also some configuration variables for network communication with Google
Cloud Storage API:

| Variable                     | Default value | Required | Description                                                                                                  |
| ---------------------------- | ------------- | -------- | ------------------------------------------------------------------------------------------------------------ |
| GCS_CLIENT_TIMEOUT           | 2s            | No       | Hard timeout on requests that gcs-helper sends to the Google Storage API                                     |
| GCS_CLIENT_IDLE_CONN_TIMEOUT | 120s          | No       | Maximum duration of idle connections between gcs-helper and the Google Storage API                           |
| GCS_CLIENT_MAX_IDLE_CONNS    | 10            | No       | Maximum number of idle connections to keep open. This doesn't control the maximum number of connections      |

### GCS_HELPER_PROXY_TIMEOUT x GCS_CLIENT_TIMEOUT

The timeout configuration is mainly controlled by two environment variables:
``GCS_HELPER_PROXY_TIMEOUT`` and ``GCS_CLIENT_TIMEOUT``. The
``GCS_HELPER_PROXY_TIMEOUT`` controls how long requests to gcs-helper can take,
and ``GCS_CLIENT_TIMEOUT`` controls how long requests from gcs-helper to
Google's API can take. Since gcs-helper automatically retries on failures, the
number of retries is roughly the value of ``GCS_HELPER_PROXY_TIMEOUT`` divided
by the value of ``GCS_CLIENT_TIMEOUT``.
