import msx from 'lib/msx';
import * as m from 'mithril';
import * as API from 'lib/api';
import Theme from 'lib/theme';
import { MustAuthController } from 'components/auth';

let _: Mithril.Component<{}, TemplatePage> = TemplatePage;

export default class TemplatePage extends MustAuthController {
  template: API.ThemeTemplate;

  constructor() {
    super();
    this.template = {};
    Theme.getFullTemplate(m.route.param('name'), m.route.param('template'))
      .then((t) => {
        this.template = t;
        m.redraw();
      });
    // todo: catch
  }

  colorize(el: HTMLElement) {
    require.ensure([
      'highlight.js',
      'highlight.js/lib/languages/xml',
      'highlight.js/styles/rainbow.css'
    ], (require) => {
      let hljs: any = require('highlight.js');
      require('highlight.js/lib/languages/xml');
      require('highlight.js/styles/rainbow.css');
      hljs.highlightBlock(el);
    }, 'hljs');
  }

  static oninit(v: Mithril.Vnode<{}, TemplatePage>) {
    v.state = new TemplatePage();
  };

  static view(v: Mithril.Vnode<{}, TemplatePage>) {
    let name = m.route.param('name');
    let ctrl = v.state;
    let lst: Mithril.Vnode<any, any>[] = [];
    let p = ctrl.template.placeholders;
    if (p && p.length > 0) {
      p.forEach((placeholder) => {
        if (placeholder.key == 'content') {
          return;
        }
        lst.push(<div class='tr'>{placeholder.key}</div>);
      });
    }
    if (!ctrl.template.hideContent) {
      lst.push(<div class='tr'>content</div>);
    }

    let placeholders: Mithril.Vnode<any, any>;
    if (lst.length > 0) {
      placeholders = <div>
        <h2>Fields</h2>
        <div class='table'>{lst}</div>
      </div>;
    }

    return <div class='template'>
      <header>
        <h1>
          <a href='/admin/themes' oncreate={m.route.link}>
            Themes
          </a>
          {m.trust(' &rsaquo; ')}
          <a href={`/admin/themes/${name}`}
            class='unbold'
            oncreate={m.route.link}
          >
            {name}
          </a>
          {m.trust(' &rsaquo; ')}
          <span class='unbold'>
            {m.route.param('template')}
          </span>
        </h1>
      </header>
      {placeholders}
      <h2>Template</h2>
      <pre onupdate={(v: Mithril.VnodeDOM<any, any>) => {
          ctrl.colorize(v.dom as HTMLElement);
        }}
      >
        {ctrl.template.data}
      </pre>
    </div>;
  }
}