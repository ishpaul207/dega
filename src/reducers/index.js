import { combineReducers } from 'redux';
import settings from './settings';
import spaces from './spaces';
import categories from './categories';
import tags from './tags';
import media from './media';

export default combineReducers({
  settings,
  spaces,
  categories,
  tags,
  media,
});
