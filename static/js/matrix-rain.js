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
        "<input type='text' onfocus=alert('Gotcha')>"
    ];

    // Configuration for Matrix Rain Effect
    const fontSize = 15;
    const lineHeight = fontSize * 4.2;
    const columns = Math.floor(canvas.width / fontSize);
    const rainDrops = Array(columns).fill(1);
    const columnScripts = Array(columns).fill(null); // Track active script for each column
    const scriptIndexes = Array(columns).fill(0); // Track character index for each active script

    function draw() {
        context.fillStyle = 'rgba(0, 0, 0, 0.05)';
        context.fillRect(0, 0, canvas.width, canvas.height);

        context.fillStyle = '#0F0';
        context.font = `${fontSize}px monospace`;

        rainDrops.forEach((y, index) => {
            let text;

            // Check if this column is currently displaying a script
            if (columnScripts[index]) {
                const script = columnScripts[index];
                const charIndex = scriptIndexes[index];

                // Get the next character in the script
                text = script[charIndex];

                // Move to the next character or reset if done
                if (charIndex < script.length - 1) {
                    scriptIndexes[index]++;
                } else {
                    columnScripts[index] = null; // Reset script for this column
                    scriptIndexes[index] = 0;
                }
            } else if (Math.random() > 0.995) {
                // Randomly start a new script in this column
                columnScripts[index] = scriptMessages[Math.floor(Math.random() * scriptMessages.length)];
                text = columnScripts[index][0]; // Start with the first character
                scriptIndexes[index] = 1; // Set the next index to 1
            } else {
                // Regular random character rain
                text = alphabet[Math.floor(Math.random() * alphabet.length)];
            }

            // Render the character
            context.fillText(text, index * fontSize, y * lineHeight);

            // Reset raindrop to the top if it reaches the bottom
            if (y * lineHeight > canvas.height && Math.random() > 0.995) {
                rainDrops[index] = 0;
                columnScripts[index] = null; // Reset script if active
                scriptIndexes[index] = 0;
            }

            // Control the speed of each raindrop
            rainDrops[index] += Math.random() * 0.4 + 0.2;
        });
    }

    // Refresh the screen every 60ms for smooth animation
    setInterval(draw, 50);
});

/* document.addEventListener("DOMContentLoaded",()=>{let e=document.getElementById("matrixRainCanvas"),t=e.getContext("2d");function l(){e.width=window.innerWidth,e.height=window.innerHeight}l(),window.addEventListener("resize",l);let n="アァィイゥウェエォオカガキギクグケゲコゴサザシジスズセゼソゾタダチヂツヅテデトドナニヌネノハバパヒビピフブプヘベペホボポマミムメモヤャユュヨョラリルレロワヲンヴヵヶABCDEFGHIJKLMNOhacktheworldPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()_+-={}[]|:;\"'<>,.?/~`",r=["<script>alert('Hacked!');</script>","<img src=x onerror=alert(1)>","<iframe src='javascript:alert(1)'></iframe>","<marquee>Hack the World!</marquee>","<script>console.log('Pwned');</script>","<input type='text' onfocus=alert('Gotcha')>"],i=63,o=Math.floor(e.width/15),a=Array(o).fill(1),d=Array(o).fill(null),h=Array(o).fill(0);function $(){t.fillStyle="rgba(0, 0, 0, 0.05)",t.fillRect(0,0,e.width,e.height),t.fillStyle="#0F0",t.font="15px monospace",a.forEach((l,i)=>{let o;if(d[i]){let $=d[i],c=h[i];o=$[c],c<$.length-1?h[i]++:(d[i]=null,h[i]=0)}else Math.random()>.995?(d[i]=r[Math.floor(Math.random()*r.length)],o=d[i][0],h[i]=1):o=n[Math.floor(Math.random()*n.length)];t.fillText(o,15*i,63*l),63*l>e.height&&Math.random()>.995&&(a[i]=0,d[i]=null,h[i]=0),a[i]+=.4*Math.random()+.2})}setInterval($,50)}); */

