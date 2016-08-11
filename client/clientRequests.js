'use strict'

const http = require('http')

module.exports = class Requester {
	constructor(hostname, port) {
		this._hostname = hostname
		this._port = port
	}

	statusAll() {
		return requestPromise({
			hostname: this._hostname,
			port: this._port,
			path: '/status',
			method: 'GET'
		})
	}

	statusOne(taskName) {
		return requestPromise({
			hostname: this._hostname,
			port: this._port,
			path: '/status/' + taskName,
			method: 'GET'
		})
	}

	startOne(taskName) {
		return requestPromise({
			hostname: this._hostname,
			port: this._port,
			path: '/start/' + taskName,
			method: 'GET'
		})
	}

	stopOne(taskName) {
		return requestPromise({
			hostname: this._hostname,
			port: this._port,
			path: '/stop/' + taskName,
			method: 'GET'
		})
	}

	restartOne(taskName) {
		return requestPromise({
			hostname: this._hostname,
			port: this._port,
			path: '/restart/' + taskName,
			method: 'GET'
		})
	}

	shutdown() {
		return requestPromise({
			hostname: this._hostname,
			port: this._port,
			path: '/shutdown',
			method: 'GET'
		})
	}
}

const requestPromise = options => {
	return new Promise((resolve, reject) => {
		const req = http.request(options, res => {
			let content = ''

			res.on('data', chunk => content += chunk)
			res.on('end', () => {
				if (res.statusCode / 100 !== 2) {
					reject(content)
				} else {
					resolve(content)
				}
			})
		})
		req.on('error', err => reject(err))
		req.end()
	})
}
