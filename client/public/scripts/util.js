let getUnits = (free, used) => {
    let out = [];
    
    if (free.length > 1 && used.length > 2) {
        let freeUnit = free.substring(free.length - 1,  free.length);
        let usedUnit = used.substring(used.length - 1,  used.length);
        out.push(freeUnit, usedUnit);
    }

    return out;
}

let convertToSame = (free, used) => {
    let out = [];
    
    if (free.length > 1 && used.length > 2) {
        let freeUnit = free.substring(free.length - 1,  free.length);
        let usedUnit = used.substring(used.length - 1,  used.length);
        let freeAmount = parseFloat(free.substring(0,  free.length - 1));
        let usedAmount = parseFloat(used.substring(0,  used.length - 1));

        if (freeUnit !== usedUnit) {
            usedAmount = convertTo(usedAmount, usedUnit, freeUnit);
        }        
        out.push(freeAmount, usedAmount);
    }

    return out;
}

let convertTo = (amount, unit, outUnit) => {
    let out = null;
    switch (unit) {
        case 'B':
            if (outUnit === 'M') {
                out = (amount / 1024) / 1024;
            } else if (outUnit === 'K') {
                out = amount / 1024;
            }
            break;
        case 'M':
            if (outUnit === 'G') {
                out = amount / 1024;
            } else if (outUnit === 'T') {
                out = (amount / 1024) / 1024;
            }
            break;
        case 'M':
            if (outUnit === 'G') {
                out = amount / 1024;
            } else if (outUnit === 'T') {
                out = (amount / 1024) / 1024;
            }
            break;
        case 'G':
            if (outUnit === 'M') {
                out = amount * 1024;
            } else if (outUnit === 'T') {
                out = amount / 1024;
            }
            break;
        case 'T':
            if (outUnit === 'M') {
                out = (amount * 1024) * 1024;
            } else if (outUnit === 'G') {
                out = amount * 1024;
            }
            break;
    }
    return Math.round(out);
}