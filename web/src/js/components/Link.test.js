import { fireEvent, render } from '@testing-library/svelte'
import sinon from 'sinon';

import Link from './Link.svelte';
import routingService from '../services/routingService';

const sandbox = sinon.createSandbox();

beforeEach(() => {
  sandbox.spy(routingService);
});

afterEach(() => {
  sandbox.restore();
});

test('Renders correctly', () => {
  const path = '/some/path';
  const rendered = render(Link, { href: path });
  const link = rendered.container.querySelector('a');

  expect(link).toHaveAttribute('href', path)
});

test('Calls goTo when clicked', async () => {
  const path = '/some/path';
  const rendered = render(Link, { href: path });
  const link = rendered.container.querySelector('a');

  await fireEvent.click(link);

  expect(routingService.goTo.calledWith(path)).toBeTruthy()
});
