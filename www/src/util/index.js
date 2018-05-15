export function hexToRgb(hex) {
    var result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
    return result ? {
        r: parseInt(result[1], 16),
        g: parseInt(result[2], 16),
        b: parseInt(result[3], 16)
    } : null;
}

import _ from 'underscore';
export const removeDuplicates = (identity: (any) => any) => {
    var seen = [];
    return (x) => {
        let id = identity(x); 
        if(_.contains(seen, id)) {
            return false;
        }
        seen.push(id)
        return true
    }
}

export function toggleInArray(arr, val) {
    if(_.contains(arr, val)) {
        return arr.filter(x => !(x == val));
    } else {
        return [...arr, val]
    }
}