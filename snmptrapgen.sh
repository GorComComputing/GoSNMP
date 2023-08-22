#!/bin/sh

PATH=$path:/bin:/usr/bin:/usr/ucb        
DIRETC=/home/icbcom/etc
FILETRAPGEN="$DIRETC/snmptrapgen.conf"
FILETRAPACTION=/home/icbcom/etc/snmpalarmaction.conf

DIR="/home/icbcom"

FILELOG="/tmp/snmptrap.log"

TRAPIP="192.168.63.68"


DEVOID=".1.3.6.1.4.1.46667.7"    


if [ -e $FILETRAPGEN ]
then
	. $FILETRAPGEN
fi



DEVOIDTRAP="$DEVOID.1"
OIDVALUE="$DEVOIDTRAP.0.1"
OIDTEXT="$DEVOIDTRAP.0.2"
OIDCRITICALY="$DEVOIDTRAP.0.3"




NAME=$(echo $1 | tr -d '"')
TRAPOID="$DEVOIDTRAP.$2" 
TEXT=$(echo $3 | tr -d '"')
CRITICALY=$4
VALUE=$(echo $5 | tr -d '"')

DATA=""
if [ "$5" ]
then
    DATA=$5
fi


#echo "NAME=$NAME"
#echo "TEXT=$TEXT"
#echo "VALUE=$VALUE"

#FILETRAPOID="/tmp/trap$(echo $2 | tr -d '.')"
#OLDVALUE="-"


#if [ ! $VALUE = $OLDVALUE ]
#then
#	if [ $VALUE = "0" ] 
#	then
#		VALUE="$NAME: END ALARM"
#		echo "OLDVALUE=0" > $FILETRAPOID
#		echo "TIME=\"$(date)\""  >> $FILETRAPOID
#		echo "$(date) $VALUE $DATA" >> $FILELOG
#	else
#		VALUE="$NAME: BEGIN ALARM"
#		echo "OLDVALUE=1" > $FILETRAPOID
#		echo "TIME=\"$(date)\""  >> $FILETRAPOID
#		echo "$(date) $VALUE $DATA" >> $FILELOG
#	fi
CMD="snmptrap -c public -v 2c -t10 $TRAPIP \"\" $TRAPOID $OIDVALUE s \"$VALUE\" $OIDTEXT s \"$TEXT\" $OIDCRITICALY i $CRITICALY"
echo $CMD
	if [ "$6" ]
	then
		ACTION=$6
	    echo "ACTION=$ACTION"
		if [ -e $FILTRAPACTION ]
	    then	
echo "START ACTION RUN"
				STR=$(grep "^$ACTION"  $FILETRAPACTION | sed -e 's/ /_/g'| sed -e 's/\\;/||/g')
echo "$STR"				
if [ "$STR" ] 
				then
	    			
echo "VALUE=$VALUE"
	    			if [ "$VALUE" = "1" ]
	    			then
						RUN=$(echo $STR|awk -F ";" '{print $2}'| sed -e 's/_/ /g' | sed -e 's/||/;/g')
	    			else
						RUN=$(echo $STR|awk -F ";" '{print $3}'| sed -e 's/_/ /g' | sed -e 's/||/;/g')
	    			fi
	    			echo "RUN=$RUN"
	    			RUN=$(eval $RUN) 
				echo $RUN
				fi
		fi
	fi

	if [ -e $FILETRAPGEN ]
	then
echo "START TRAP SENDING"	    
TRAPIPS=$(grep "TRAPIP" $FILETRAPGEN|awk -F "=" '{print $2}' )
	    for TRAPIP in $TRAPIPS
	    do
			echo "TRAPIP=$TRAPIP"
		if [ "$TRAPIP" != "none" ]
		then
			LOP=$(ping -c 3 -w 5 $TRAPIP 2> /dev/null | grep 'packet loss'| awk -F " " '{print $7}'| awk -F "%" '{print $1}')
			#$CMD
			if [ "$LOP" == "0" ]
			then
				snmptrap -c public -v 2c -t10 $TRAPIP "" $TRAPOID $OIDVALUE s "$VALUE" $OIDTEXT s "$TEXT" $OIDCRITICALY i $CRITICALY 
			fi
		fi
		#snmptrap $CMD
	    done
	fi
#fi

echo "END SCRIPT"

exit 0



