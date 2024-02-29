import chalk from 'chalk';
import chokidar from 'chokidar';
import { copyFileSync, existsSync, lstatSync, mkdirSync, readFileSync, rmdirSync, unlink, unlinkSync } from 'fs';
import ignore from 'ignore';
import { basename, dirname, sep } from 'path';
import ora from 'ora';

/**
 * Prints error message and exits
 * @param {String} text
 */
const error = (text) => {
	console.log(chalk.redBright(text));
	process.exit();
}

/**
 * Checks to see if path is valid
 * @param {String} path
 * @returns Boolean
 */
const isDirectory = (path) => {
	let isDir = false;
	try {
		isDir = lstatSync(path).isDirectory();
	} catch (e) {
	}
	return isDir;
}

/**
 * Prints syntax for script
 * @param {Boolean} exit Exit after displaying syntax?
 */
const syntax = (exit = true) => {
	console.log(chalk.greenBright('Syntax: npm start source-directory destination-directory\n'));
	if (process.platform === 'win32') {
		console.log(chalk.whiteBright('Example: npm start C:\\code\\my-app \\\\server-name\\code-share\n'));
	} else {
		console.log(chalk.whiteBright('Example: npm start ~/code/my-app /Volumes/mount/network-share\n'));
	}
	if (exit) process.exit();
}

/**
 * Copies file or directory to destination
 * @param {String} src Path of file or dir to sync
 * @param {String} action 'copy' or 'delete'
 */
const syncPath = async (src, action) => {
	const relativePath = src.replace(dir2watch, '');
	const dst = src.replace(dir2watch, destdir);
	const dstDir = dirname(dst);

	if (ig.ignores(relativePath)) {
		console.log(chalk.cyanBright(`Ignoring ${relativePath}`));
		return;
	}

	if (action === 'copy') {
		const spinner = ora({
			text: `Updating ${relativePath}...`,
			color: 'green',
			spinner: 'circleHalves',
			interval: 200,
		}).start();
		try {
			if (isDirectory(src)) {
				mkdirSync(dst, { recursive: true});
			} else {
				mkdirSync(dstDir, { recursive: true});
				await copyFileSync(src, dst);
			}
			spinner.succeed(`${relativePath} copied ${new Date()}`);
			// Remove hidden files created by Mac
			const hiddenFile = `${dstDir}${sep}._${basename(dst)}`;
			setTimeout(() => {
				unlink(hiddenFile, () => {})
			}, 1000);
		} catch (e) {
			spinner.fail(`${relativePath} failed to copy`);
			console.log(e)
		}
	}
	if (action === 'delete') {
		const spinner = ora({
			text: `Removing ${relativePath}...`,
			color: 'green',
			spinner: 'circleHalves',
			interval: 200,
		}).start();
		try {
			if (isDirectory(dst)) {
				await rmdirSync(dst);
			} else {
				await unlinkSync(dst);
			}
			spinner.succeed(`${relativePath} removed ${new Date()}`);
		} catch (e) {
			spinner.fail(`Failed to remove ${relativePath}`);
			console.log(e)
		}
	}
}

// Check for valid arguments.  Print syntax if necessary.
if (process.argv.length < 4) {
	syntax();
}

// Read arguments and verify
const parsePath = (path) => {
	const lastChar = path.charAt(path.length - 1);
	if (lastChar === sep) return path;
	return path + sep;
}
const dir2watch = parsePath(process.argv[2]);
const destdir = parsePath(process.argv[3]);

if (isDirectory(dir2watch) === false) {
	syntax(false);
	error(`${dir2watch} not found or is not a directory`);
}
if (isDirectory(destdir) === false) {
	syntax(false);
	error(`${destdir} not found or is not a directory`);
}


// Start watching for changes
const watcher = chokidar.watch(dir2watch, {
	ignored: /(^|[\/\\])\../, // ignore dotfiles
	persistent: true,
	ignoreInitial: true
});

console.log(chalk.yellowBright(`Monitoring ${dir2watch} for changes\n`));

const ig = ignore();
const ignoreFile = `${dir2watch}/.csignore`;
if (existsSync(ignoreFile)) {
	const fileData = readFileSync(ignoreFile).toString();
	ig.add(fileData.split(/\n/));
}

watcher
	.on('add', (path) => syncPath(path, 'copy'))
	.on('addDir', (path) => syncPath(path, 'copy'))
	.on('change', (path) => syncPath(path, 'copy'))
	.on('unlink', (path) => syncPath(path, 'delete'))
	.on('unlinkDir', (path) => syncPath(path, 'delete'))
	.on('error', (error) => console.error(chalk.redBright('An error occurred:'), error));
