import hljs from 'highlight.js/lib/common';
import dockerfile from 'highlight.js/lib/languages/dockerfile';
import scala from 'highlight.js/lib/languages/scala';

hljs.registerLanguage('dockerfile', dockerfile);
hljs.registerLanguage('scala', scala);

export default hljs;
