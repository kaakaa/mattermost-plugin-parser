const winston = require('winston');

const logFormat = winston.format.printf(({ level, message, label, timestamp}) => {
    return `${timestamp} [${label}] ${level}: ${message}`;
});

const logger = winston.createLogger({
    format: winston.format.combine(
        winston.format.label({label: 'mmparser'}),
        winston.format.timestamp(),
        winston.format.splat(),
        winston.format.colorize(),
        logFormat,
    ),
    transports: [
        new winston.transports.Console(),
    ]
});

module.exports = {
    logger: logger,
};