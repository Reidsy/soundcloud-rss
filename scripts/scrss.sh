# Used to add start, stop, restart commands to your shell.
# Usage: scrss {start|stop|restart}
#
# Include this in a .profile by adding `source scrss.sh`
#
# The script requires `soundcloud-rss.env` to exist in the current directory.

scrss () {
  case "$1" in
    start)
      echo "Starting soundcloud-rss..."
      docker pull docker.pkg.github.com/reidsy/soundcloud-rss/master
      docker run -d -p 80:8080 --env-file ./soundcloud-rss.env --name soundcloud-rss docker.pkg.github.com/reidsy/soundcloud-rss/master
      ;;
    stop)
      echo "Stopping soundcloud-rss..."
      docker stop soundcloud-rss
      docker rm soundcloud-rss
      ;;
    restart)
      echo "Stopping soundcloud-rss..."
      docker stop soundcloud-rss
      docker rm soundcloud-rss
      echo "Starting soundcloud-rss..."
      docker pull docker.pkg.github.com/reidsy/soundcloud-rss/master
      docker run -d -p 80:8080 --env-file ./soundcloud-rss.env --name soundcloud-rss docker.pkg.github.com/reidsy/soundcloud-rss/master
      ;;
    *)
      echo "Usage: scrss {start|stop|restart}"
  esac
}
