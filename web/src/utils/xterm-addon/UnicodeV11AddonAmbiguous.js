import { UnicodeV11Ambiguous } from './UnicodeV11Ambiguous';

export class UnicodeV11AddonAmbiguous {
  activate(terminal) {
    terminal.unicode.register(new UnicodeV11Ambiguous());
  }
  dispose() { }
}
