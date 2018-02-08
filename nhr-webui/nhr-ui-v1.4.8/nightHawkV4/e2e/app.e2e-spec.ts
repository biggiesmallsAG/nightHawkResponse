import { NightHawkV4Page } from './app.po';

describe('night-hawk-v4 App', function() {
  let page: NightHawkV4Page;

  beforeEach(() => {
    page = new NightHawkV4Page();
  });

  it('should display message saying app works', () => {
    page.navigateTo();
    expect(page.getParagraphText()).toEqual('app works!');
  });
});
