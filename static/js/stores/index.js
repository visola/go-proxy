import Mappings from './Mappings';
import PossibleValues from './PossibleValues';
import ProxiedRequests from './ProxiedRequests';
import SelectedValues from './SelectedValues';
import Variables from './Variables';

const mappings = new Mappings();
mappings.fetch();

const possibleValues = new PossibleValues();
possibleValues.fetch();

const proxiedRequests = new ProxiedRequests();

const selectedValues = new SelectedValues();
selectedValues.fetch();

const variables = new Variables();
variables.fetch();

export default {
  mappings,
  possibleValues,
  proxiedRequests,
  selectedValues,
  variables,
};
