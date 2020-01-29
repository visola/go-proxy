import sinon from 'sinon';

import routingService from '../services/routingService';

const sandbox = sinon.createSandbox();

beforeEach(() => {
  sandbox.spy(history);
});

afterEach(() => {
  sandbox.restore();
});

test('Calls history.pushState when in goTo', () => {
  const testPath = '/some/path';
  routingService.goTo(testPath);

  expect(history.pushState.calledWith(null, null, testPath)).toBeTruthy()

  const value = routingService.getValue();
  expect(value).toBe(testPath);
});
