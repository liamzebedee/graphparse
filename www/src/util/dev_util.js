import { API_BACKEND_HOST } from './env/config';

let lastMod = null;
function listenForReload(cb) {
    let poll;

    poll = () => {
        fetch(`${API_BACKEND_HOST}/graph/last-generated`)
        .then(currentMod => {
            if(lastMod < currentMod) {
                cb()
            }
        })
        
        setTimeout(poll, 1000)
    }
    
    poll()
}

export { listenForReload };