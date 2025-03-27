document.addEventListener('DOMContentLoaded', () => {
    const canvas = document.getElementById('matrixRainCanvas');
    const context = canvas.getContext('2d');

    // Resize the canvas to fit the screen
    function resizeCanvas() {
        canvas.width = window.innerWidth;
        canvas.height = window.innerHeight;
    }
    resizeCanvas();
    window.addEventListener('resize', resizeCanvas);

    // Full Character Set for Matrix Rain
    const katakana =
        'アァィイゥウェエォオカガキギクグケゲコゴサザシジスズセゼソゾタダチヂツヅテデトドナニヌネノハバパヒビピフブプヘベペホボポマミムメモヤャユュヨョラリルレロワヲンヴヵヶ';
    const latin = 'ABCDEFGHIJKLMNOhacktheworldPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz';
    const nums = '0123456789';
    const symbols = '!@#$%^&*()_+-={}[]|:;"\'<>,.?/~`';

    const alphabet = katakana + latin + nums + symbols;

    // Array of Specific Script Messages
    const scriptMessages = [
        "<script>alert('Hacked!');</script>",
        "<img src=x onerror=alert(1)>",
        "<iframe src='javascript:alert(1)'></iframe>",
        "<marquee>Hack the World!</marquee>",
        "<script>console.log('Pwned');</script>",
        "<input type='text' onfocus=alert('Gotcha')>",
        "Wake up, Neo...",
        "The Matrix has you...",
        "Follow the white rabbit",
        "Knock, knock, Neo",
        "There is no spoon"
    ];

    // Configuration for Matrix Rain Effect
    const fontSize = 15;
    // Double the spacing between characters from the previous 1.5x
    const charSpacing = fontSize * 3.0; // Even wider spacing between characters
    const columns = Math.floor(canvas.width / charSpacing);
    const rainDrops = Array(columns).fill(1);
    const columnScripts = Array(columns).fill(null);
    const scriptIndexes = Array(columns).fill(0);
    const speeds = Array(columns).fill(0).map(() => Math.random() * 0.5 + 0.3);
    
    // Reserved columns for script messages to prevent overlap
    const reservedColumns = new Set();

    function draw() {
        context.fillStyle = 'rgba(0, 0, 0, 0.05)';
        context.fillRect(0, 0, canvas.width, canvas.height);

        for (let i = 0; i < columns; i++) {
            const y = rainDrops[i] * fontSize * 1.2;
            // Use wider horizontal spacing
            const x = i * charSpacing;
            
            let text;
            
            // Handle script messages with better spacing
            if (columnScripts[i]) {
                const script = columnScripts[i];
                const charIndex = scriptIndexes[i];
                text = script[charIndex];
                
                // Lead character is bright white
                context.fillStyle = '#FFFFFF';
                
                if (charIndex < script.length - 1) {
                    scriptIndexes[i]++;
                } else {
                    columnScripts[i] = null;
                    scriptIndexes[i] = 0;
                    reservedColumns.delete(i);
                }
            } else if (Math.random() > 0.995 && !reservedColumns.has(i)) {
                // Ensure there's space for the message
                let hasSpace = true;
                const freeColumns = 5; // Keep space around message columns
                
                for (let j = Math.max(0, i - freeColumns); j <= Math.min(columns - 1, i + freeColumns); j++) {
                    if (j !== i && reservedColumns.has(j)) {
                        hasSpace = false;
                        break;
                    }
                }
                
                if (hasSpace) {
                    columnScripts[i] = scriptMessages[Math.floor(Math.random() * scriptMessages.length)];
                    text = columnScripts[i][0];
                    scriptIndexes[i] = 1;
                    reservedColumns.add(i);
                    context.fillStyle = '#FFFFFF';
                } else {
                    text = alphabet[Math.floor(Math.random() * alphabet.length)];
                    
                    // Regular characters with varied green intensity
                    const brightness = 160 + Math.floor(Math.random() * 95);
                    context.fillStyle = `rgb(0, ${brightness}, 0)`;
                }
            } else {
                text = alphabet[Math.floor(Math.random() * alphabet.length)];
                
                // First character in stream is brighter
                if (rainDrops[i] <= 1) {
                    context.fillStyle = '#FFFFFF';
                } else {
                    const brightness = 160 + Math.floor(Math.random() * 95);
                    context.fillStyle = `rgb(0, ${brightness}, 0)`;
                }
            }
            
            context.font = `${fontSize}px monospace`;
            context.fillText(text, x, y);
            
            // Reset raindrop when it reaches bottom
            if (y > canvas.height && Math.random() > 0.975) {
                rainDrops[i] = 0;
                speeds[i] = Math.random() * 0.5 + 0.3;
                
                if (columnScripts[i]) {
                    columnScripts[i] = null;
                    scriptIndexes[i] = 0;
                    reservedColumns.delete(i);
                }
            }
            
            // Update position with consistent but varied speed
            rainDrops[i] += speeds[i];
        }
    }

    // Faster refresh rate for smoother animation
    setInterval(draw, 30);
});

/* document.addEventListener("DOMContentLoaded",()=>{let e=document.getElementById("matrixRainCanvas"),t=e.getContext("2d");function l(){e.width=window.innerWidth,e.height=window.innerHeight}l(),window.addEventListener("resize",l);let n="アァィイゥウェエォオカガキギクグケゲコゴサザシジスズセゼソゾタダチヂツヅテデトドナニヌネノハバパヒビピフブプヘベペホボポマミムメモヤャユュヨョラリルレロワヲンヴヵヶABCDEFGHIJKLMNOhacktheworldPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()_+-={}[]|:;\"'<>,.?/~`",r=["<script>alert('Hacked!');</script>","<img src=x onerror=alert(1)>","<iframe src='javascript:alert(1)'></iframe>","<marquee>Hack the World!</marquee>","<script>console.log('Pwned');</script>","<input type='text' onfocus=alert('Gotcha')>"],i=63,o=Math.floor(e.width/15),a=Array(o).fill(1),d=Array(o).fill(null),h=Array(o).fill(0);function $(){t.fillStyle="rgba(0, 0, 0, 0.05)",t.fillRect(0,0,e.width,e.height),t.fillStyle="#0F0",t.font="15px monospace",a.forEach((l,i)=>{let o;if(d[i]){let $=d[i],c=h[i];o=$[c],c<$.length-1?h[i]++:(d[i]=null,h[i]=0)}else Math.random()>.995?(d[i]=r[Math.floor(Math.random()*r.length)],o=d[i][0],h[i]=1):o=n[Math.floor(Math.random()*n.length)];t.fillText(o,15*i,63*l),63*l>e.height&&Math.random()>.995&&(a[i]=0,d[i]=null,h[i]=0),a[i]+=.4*Math.random()+.2})}setInterval($,50)}); */

