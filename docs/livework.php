<?php
	// ******************************************************
	//    Livetiming mit EasyWk 5.xx
	//		Diese Datei wird vom EasyWk angesprochen und
	//		erzeugt die eigentlichen Livetiming-Daten
	//		Sie muss im Livetiming-Ordner liegen, wird
	//		aber nicht verklinkt!!!
	//	  Stand: 05.04.2022
	// ******************************************************

	// =====================================
	// Initialisierung
	// =====================================
	header('Content-Type: text/html; charset=utf-8');
	error_reporting(E_ALL);
	require_once('livebasic.php');
	$dolivedata = FALSE;
	$dowritefile = FALSE;
	$doracesummary = FALSE;
	$result = '';

	// =====================================
	// Schwimmzeit als Text
	// =====================================
	function swimtime2str($swimtime) {
		if ($swimtime>0) {
			$swimtime=str_pad($swimtime,8,"0",STR_PAD_LEFT);
			if (substr($swimtime,0,2)=='00')
			  return substr($swimtime,2,2).':'.substr($swimtime,4,2).','.substr($swimtime,6,2);
			else return substr($swimtime,0,2).':'.substr($swimtime,2,2).':'.substr($swimtime,4,2).','.substr($swimtime,6,2);
		} else {
			return '&nbsp;';
		}
	}

	// =====================================
	// Reaktionszeit als Text
	// =====================================
	function rt2str($swimtime) {
		if ($swimtime>0) {
			$swimtime=str_pad($swimtime,3,"0",STR_PAD_LEFT);
			return substr($swimtime,0,1).','.substr($swimtime,1,2);
		} else {
			return '&nbsp;';
		}
	}

	// =====================================
	// Hauptteil
	// =====================================
	// Sicherheitsabfrage
	if (isset($_REQUEST['pwd']) && ($_REQUEST['pwd']==defLIVEPWD) && isset($_REQUEST['action'])) {

		// Datei einbinden wenn vorhanden, ansonsten Array intialisieren
		if (file_exists('data.php')) {
			if (function_exists('opcache_invalidate')) { opcache_invalidate('data.php',true); }
			include_once('data.php');
			if (!isset($base['event']))
				$base['event'] = 0;
			if (!isset($base['name']))
				$base['name'] = '';
		} else {
			$base = array();
			$base['firstlane'] = 1;
			$base['lanecount'] = 10;
			$base['vername'] = '';
			$base['lastevent'] = 0;
			$base['event'] = 0;
			$base['heat'] = 0;
			$base['maxheat'] = 0;
			$base['name'] = '';
			$lanes = array();
			$dowritefile = TRUE;
		}

		// Action gibt an, was gesendet wird
		switch ($_REQUEST['action']) {

			case 'ping':
				$result = 'OK';
				break;

			case 'clearsum':
				$summary = array();
				file_put_contents('summary.php',"<?php\n"
					.'$summary = '.var_export($summary,true).";\n"
					."?>\n");
				$result = 'OK';
				$text = '<p>&nbsp;</p>'."\n";
				file_put_contents('livesum.php',$text);
				$base['lastevent'] = 0;
				$dowritefile = TRUE;
				break;

			case 'init':
				// eventuell soll das Rennen beibehalten werden
				// da Daten "hinterher geschoben" werden
				if (isset($_REQUEST['keepsum']))
					$lastevent = $base['lastevent'];
				else $lastevent = 0;
				// Eingetragene Daten löschen
				unset($base); $base = array();
				unset($lanes); $lanes = array();
				// Grunddaten übernehmen
				$base['firstlane']=intval($_REQUEST['firstlane']);
				$base['lanecount']=intval($_REQUEST['lanecount']);
				$base['vername']=utf8_encode(str_replace(array('<?','?>'),array('',''),$_REQUEST['vername']));
				$base['lastevent']=$lastevent;
				$base['event'] = 0;
				$base['heat'] = 0;
				$base['maxheat'] = 0;
				$base['name'] = '';
				// Lanes wieder aufbauen
				for ($i=$base['firstlane']; $i<=($base['firstlane']+$base['lanecount']-1); $i++) {
					$lanes[$i]['lane'] = $i;
					$lanes[$i]['name'] = '';
					$lanes[$i]['yob'] = '';
					$lanes[$i]['club'] = '';
					$lanes[$i]['meter'] = '';
					$lanes[$i]['time'] = 0;
					$lanes[$i]['m'] = 0;
				}
				$dowritefile = TRUE;
				$result = 'OK';
				break;

			case 'newrace':
			case 'ready':
				// Neues Rennen eintragen
				$base['event'] = intval($_REQUEST['event']);
				$base['heat'] = intval($_REQUEST['heat']);
				$base['maxheat'] = intval($_REQUEST['maxheat']);
				$base['name'] = str_replace(array('<?','?>'),array('',''),utf8_encode($_REQUEST['name']));
				// Datei mit den aktuellen Zeiten leeren, wenn vorhanden Schwimmer eintragen
				for ($i=$base['firstlane']; $i<=($base['firstlane']+$base['lanecount']-1); $i++) {
					$lanes[$i]['lane'] = $i;
					if (isset($_REQUEST['swr'.$i]))
					  $lanes[$i]['name'] = str_replace(array('<?','?>'),array('',''),utf8_encode($_REQUEST['swr'.$i]));
					else $lanes[$i]['name'] = '';
					if (isset($_REQUEST['yob'.$i]))
					  $lanes[$i]['yob'] = $_REQUEST['yob'.$i];
					else $lanes[$i]['yob'] = '';
					if (isset($_REQUEST['club'.$i]))
					  $lanes[$i]['club'] = str_replace(array('<?','?>'),array('',''),utf8_encode($_REQUEST['club'.$i]));
					else $lanes[$i]['club'] = '';
					$lanes[$i]['meter'] = '';
					$lanes[$i]['time'] = 0;
					$lanes[$i]['rank'] = '';
					$lanes[$i]['m'] = 0;
				}
				$dowritefile = TRUE;
				$dolivedata = TRUE;
				$result = 'OK';
				break;

			case 'time':
				// Zeiten für eine Bahn eintragen
				$lane = intval($_REQUEST['lane']);
				$lanes[$lane]['meter'] = trim($_REQUEST['meter'])=='RT' ? 'RT' : intval($_REQUEST['meter']);
				$lanes[$lane]['time'] = intval($_REQUEST['time']);
				$lanes[$lane]['m'] = intval(str_replace(array('RT','m'),array('',''),$_REQUEST['meter']));
				// ist es eine Endzeit?
				if (isset($_REQUEST['finished']) && ($_REQUEST['finished']=='yes') && ($lanes[$lane]['time']>0)) {
					// Bestehende Daten laden
					if (file_exists('summary.php')) {
						if (function_exists('opcache_invalidate')) { opcache_invalidate('summary.php',true); }
						include_once('summary.php');
					}
					else $summary = array();
					// Sind wir im nächsten Wettkampf? dann bestehendes Summary löschen
					if ($base['event']!=$base['lastevent']) {
						unset($summary);
						$summary = array();
						$base['lastevent'] = $base['event'];
					}
					// ins Summary eintragen
					$swr = '<td>'.$lanes[$lane]['name'].'</td><td style="text-align: center;">'.$lanes[$lane]['yob'].'</td><td>'.$lanes[$lane]['club'].'</td>';
					$summary[$swr] = $_REQUEST['time'];
					asort($summary,SORT_NUMERIC);
					// Summary merken
					file_put_contents('summary.php',"<?php\n"
						.'$summary = '.var_export($summary,true).";\n"
						."?>\n");
					// Summary-Datei für Web erzeugen
					$text = '<h2>Wk '.$base['event'].' - '.$base['name'].' - Ergebnisse</h2>'."\n"
						  . '<table cellpadding="0" cellspacing="0" border="0">'."\n"
						  . '<tr><th>Schwimmer</th><th>Jg</th><th>Verein</th><th>Endzeit</th></tr>'."\n";
					$style = 'odd';
					foreach ($summary as $key => $value) {
						$style=='even' ? $style='odd' : $style='even';
						$text .= '<tr class="'.$style.'">'.$key.'<td style="text-align: right;">'.swimtime2str($value)."</td></tr>\n";
					}
					$text .= '</table>'."\n";
					// Daten in die Datei fürs Summary schreiben ....
					file_put_contents('livesum.php',$text);
				}
				// Ranking rechnen
				$laneorder = array();
				for ($i=$base['firstlane']; $i<=($base['firstlane']+$base['lanecount']-1); $i++) {
					$laneorder[$i] = $lanes[$i];
				}
				array_multisort(array_column($laneorder,'m'),SORT_DESC,array_column($laneorder,'time'),SORT_ASC,$laneorder);
				$place = 0;
				$lasttime = 0;
				$lastmeter = 0;
				foreach ($laneorder as $i => $v) {
					if (intval($v['time'])>0) {
						if (($v['time']<>$lasttime) or ($v['m']<>$lastmeter))
							$place++;
						$lasttime = $v['time'];
						$lastmeter = $v['m'];
						$lanes[$v['lane']]['rank'] = $place;
					} else {
						$lanes[$v['lane']]['rank'] = '';
					}
				}
				unset($laneorder);
				$dowritefile = TRUE;
				$dolivedata = TRUE;
				$result = 'OK';
				break;

			case 'disq':
				// Zeiten für eine Bahn eintragen
				$lane = intval($_REQUEST['lane']);
				$lanes[$lane]['time'] = 595998;
				// Bestehende Daten laden
				if (file_exists('summary.php')) {
					if (function_exists('opcache_invalidate')) { opcache_invalidate('summary.php',true); }
					include_once('summary.php');
				}
				else $summary = array();
				// ins Summary eintragen
				$swr = '<td>'.$lanes[$lane]['name'].'</td><td style="text-align: center;">'.$lanes[$lane]['yob'].'</td><td>'.$lanes[$lane]['club'].'</td>';
				$summary[$swr] = 595998;
				asort($summary,SORT_NUMERIC);
				// Summary merken
				file_put_contents('summary.php',"<?php\n"
					.'$summary = '.var_export($summary,true).";\n"
					."?>\n");
				// Summary-Datei für Web erzeugen
				$text = '<h2>Wk '.$base['event'].' - '.$base['name'].' - Ergebnisse</h2>'."\n"
					  . '<table cellpadding="0" cellspacing="0" border="0">'."\n"
					  . '<tr><th>Schwimmer</th><th>Jg</th><th>Verein</th><th>Endzeit</th></tr>'."\n";
				$style = 'odd';
				foreach ($summary as $key => $value) {
					$style=='even' ? $style='odd' : $style='even';
					$text .= '<tr class="'.$style.'">'.$key.'<td style="text-align: right;">'.($value==595998 ? 'DSQ' : swimtime2str($value))."</td></tr>\n";
				}
				$text .= '</table>'."\n";
				// Daten in die Datei fürs Summary schreiben ....
				file_put_contents('livesum.php',$text);

				// Ranking rechnen
				$laneorder = array();
				for ($i=$base['firstlane']; $i<=($base['firstlane']+$base['lanecount']-1); $i++) {
					$laneorder[$i] = $lanes[$i];
				}
				array_multisort(array_column($laneorder,'m'),SORT_DESC,array_column($laneorder,'time'),SORT_ASC,$laneorder);
				$place = 0;
				$lasttime = 0;
				$lastmeter = 0;
				foreach ($laneorder as $i => $v) {
					if ((intval($v['time'])>0) && (intval($v['time'])<>595998)) {
						if (($v['time']<>$lasttime) or ($v['m']<>$lastmeter))
							$place++;
						$lasttime = $v['time'];
						$lastmeter = $v['m'];
						$lanes[$v['lane']]['rank'] = $place;
					} else {
						$lanes[$v['lane']]['rank'] = '';
					}
				}
				unset($laneorder);
				$dowritefile = TRUE;
				$dolivedata = TRUE;
				$result = 'OK';
				break;

			case 'text':
				// freien Text direkt eintragen
				$content = $_REQUEST['content'];
				if ((strpos($content,'<?') !== false) || (strpos(strtolower($content),'<script') !== false)) {
					$result = 'ERROR: Unerlaubte Zeichen im Text';
					file_put_contents('livedata.php','');
				} else {
					// Daten in die Datei für das Livetiming schreiben
					file_put_contents('livedata.php',utf8_encode($content));
					$result = 'OK';
				}
				break;

			case 'raceresult';
				// Nach Zeiten sortieren
				$laneorder = array();
				for ($i=$base['firstlane']; $i<=($base['firstlane']+$base['lanecount']-1); $i++)
					$laneorder[$i] = $lanes[$i]['time'];
				asort($laneorder,SORT_NUMERIC);
				// Livetiming-Datei neu schreiben
				$live = '<h2>Wk '.$base['event'].' - '.$base['name'].' - Lauf '.$base['heat'].'/'.$base['maxheat']."</h2>\n"
					  . '<table cellpadding="0" cellspacing="0" border="0">'."\n"
					  . '<tr><th>&nbsp;</th><th>Schwimmer</th><th>Jg</th><th>Verein</th><th>Meter</th><th>Zeit</th></tr>'."\n";
				$place = 0;
				$lasttime = 0;
				foreach ($laneorder as $i => $time) {
					if ($time>0) {
						if ($time<>$lasttime)
							$place++;
						$live .= '<tr class="';
						($place&1)? $live.='odd':$live.='even';
						$live .= '"><td style="text-align: right;">'.($lanes[$i]['time']==595998 ? '' : $place).'.</td><td>'.$lanes[$i]['name'].'</td><td style="text-align: center;">'.$lanes[$i]['yob'].'</td><td>'.$lanes[$i]['club']
							  .  '</td><td style="text-align: right;">'.$lanes[$i]['meter'].'</td><td style="text-align: right;">';
						if ($lanes[$i]['time']==595998)
							$live.='DSQ';
						elseif ($lanes[$i]['meter']=='RT')
							$live.=rt2str($lanes[$i]['time']);
						else $live.=swimtime2str($lanes[$i]['time']);
						$live .= '</td></tr>'."\n";
						$lasttime = $time;
					}
				}
				$live .= '</table>'."\n";
				unset($laneorder);
				// Für leere Bahnen eine korrekte HTML-Darstellung gewährleisten
				$live = str_replace('></td>','>&nbsp;</td>',$live);
				// Daten in die Datei für das Livetiming schreiben
				file_put_contents('livedata.php',$live);
				$result = 'OK';
				break;

			default:
			    $result = 'ERROR: Unbekannte Aktion';
				break;
		}

		// Datei neu schreiben ...
		if ($dowritefile) {
			file_put_contents('data.php',"<?php\n"
				.'$base = '.var_export($base,true).";\n"
				.'$lanes = '.var_export($lanes,true).";\n"
				."?>\n");
		}

		// Livetiming erzeugen ...
		if ($dolivedata) {
			// Überschrift
			$live = '<h2>Wk '.$base['event'].' - '.$base['name'].' - Lauf '.$base['heat'].'/'.$base['maxheat']."</h2>\n"
				  . '<table cellpadding="0" cellspacing="0" border="0">'."\n"
				  . '<tr><th>Bahn</th><th>Schwimmer</th><th>Jg</th><th>Verein</th><th>Meter</th><th>Zeit</th><th>Platz</th></tr>'."\n";
			// die einzelnen Bahnen
			for ($i=$base['firstlane']; $i<=($base['firstlane']+$base['lanecount']-1); $i++) {
				$live .= '<tr class="';
				($i&1)? $live.='odd':$live.='even';
				$live .= '"><td style="text-align: right;">'.$lanes[$i]['lane'].'</td><td>'.$lanes[$i]['name'].'</td><td style="text-align: center;">'.$lanes[$i]['yob'].'</td><td>'.$lanes[$i]['club']
					  .  '</td><td style="text-align: right;">'.$lanes[$i]['meter'].'</td><td style="text-align: right;">';
				if ($lanes[$i]['time']==595998)
					$live.='DSQ';
				elseif ($lanes[$i]['meter']=='RT')
					$live.=rt2str($lanes[$i]['time']);
				else $live.=swimtime2str($lanes[$i]['time']);
				$live .= '</td><td style="text-align: center;">'.$lanes[$i]['rank'].'</td></tr>'."\n";
			}
			$live .= '</table>'."\n";
			// Für leere Bahnen eine korrekte HTML-Darstellung gewährleisten
			$live = str_replace('></td>','>&nbsp;</td>',$live);
			// Daten in die Datei für das Livetiming schreiben
			file_put_contents('livedata.php',$live);
		}

	} else {
		$result = 'ERROR: Passwort nicht korrekt oder keine Aktion definiert';
	}

	echo $result;
?>
