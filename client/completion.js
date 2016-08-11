'use strict'
var completionsCmd = ['status', 'load', 'reload', 'stop', 'shutdown', 'exit']
var completionsArg = []

function completer(line) {
	//console.log('line is : ||' + line + '||')
	//let tmp = line.split(/\s*/) seems not to work ?
	let tmp = line.split(' ')
	//let tmp = line.replace(/^[\s\uFEFF\xA0]+/g, '').split(/\s*/)
	//console.log(tmp)

	if (tmp.length <= 1) {
		return [
			completionsCmd.filter((c) => { return c.indexOf(line) == 0}),
			line
		]
	} else {
		return [
			completionsArg
				.filter(c => c.indexOf(tmp[tmp.length - 1]) == 0)
				.map(word => popLastWord(line) + ' ' + word),
			line
		]
	}
}

function popLastWord(sentence) {
	let n = sentence.lastIndexOf(' ')
	if (n === -1) {
		return sentence
	} else {
		return sentence.split('').slice(0, n).join('')
	}
}

function delCompletionAll() {
	completionsArg = []
}

function addCompletionKey(key) {
	if (completionsArg.indexOf(key) == -1) {
		completionsArg.push(key)
	}
}

module.exports.completer = completer
module.exports.addCompletionKey = addCompletionKey
module.exports.delCompletionAll = delCompletionAll
