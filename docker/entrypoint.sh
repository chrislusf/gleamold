#!/bin/sh

case "$1" in

  'master')
  	exec /usr/bin/gleamold $@
	;;

  'agent')
  	ARGS="--host=`hostname -i`  --dir=/data"
  	exec /usr/bin/gleamold $@ $ARGS
	;;

  *)
  	exec $@
	;;
esac
