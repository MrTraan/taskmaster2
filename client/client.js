'use strict'

const readline 		= require('readline')
const Requester		= require('./clientRequests')
const comp			= require('./completion')

const rl = readline.createInterface({
	input:		process.stdin,
	output:		process.stdout,
	completer:	comp.completer
})

const requester = new Requester('localhost', 8080)
rl.setPrompt("tmclient> ")

rl.prompt()
rl.on('line', (line) => {
	const av = line.split(' ')
	let request = undefined

	switch(av[0]) {
		case 'status':
			if (!av[1]) {
				request = 'statusAll'
			} else {
				request = 'statusOne'
			}
			break
		case 'start':
			request = 'startOne'
			break
		case 'restart':
			request = 'restartOne'
			break
		case 'stop':
			request = 'stopOne'
			break
		case 'shutdown':
			request = 'shutdown'
			break
		case 'exit':
			rl.close()
			break
		case 'help':
			helpCmd()
			rl.prompt()
			break
		default:
			console.log('unknown command')
			rl.prompt()
	}
	if (request) {
		requester[request](av[1])
		.then((response) => {
			console.log(response)
			rl.prompt()
		})
		.catch((err) => {
			console.log('Error: ' + err)
			rl.prompt()
		})
	}
})

rl.on('close', () => {
	console.log('Exiting tm client')
	process.exit(1)
})

function helpCmd() {
	console.log('List of commands:')
	console.log('status	[proc]	-> get the process informations')
	console.log('start [proc]	-> start a processus')
	console.log('stop [proc]	-> stop a processus')
	console.log('restop [proc]	-> stop a processus')
	console.log('exit			-> exit the client')
	console.log('shutdown		-> stop the server')
}
