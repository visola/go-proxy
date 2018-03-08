import PropTypes from 'prop-types';
import React from 'react';

export default class Modal extends React.Component {
  static propTypes = {
    children: PropTypes.oneOfType([
      PropTypes.arrayOf(PropTypes.node),
      PropTypes.node
    ]).isRequired,
    onClose: PropTypes.func.isRequired,
    title: PropTypes.string.isRequired,
  }

  constructor(props) {
    super(props);
    this.handleKeyDown = this.handleKeyDown.bind(this);
  }

  componentDidMount(){
    document.addEventListener("keydown", this.handleKeyDown, false);
  }

  componentWillUnmount(){
    document.removeEventListener("keydown", this.handleKeyDown, false);
  }

  handleKeyDown(e) {
    if (e.keyCode === 27) {
      this.props.onClose();
    }
  }

  render () {
    return <div
      className="modal in"
      ref={(el) => el ? el.focus() : ''}
      role="dialog"
      style={{display: 'block' }}
      tabIndex="-1"
    >
      <div className="modal-dialog" role="document">
        <div className="modal-content">
          <div className="modal-header">
            <button
              type="button"
              className="close"
              aria-label="Close"
              onClick={this.props.onClose}>
              <span aria-hidden="true">&times;</span>
            </button>
            <h4 className="modal-title">{this.props.title}</h4>
          </div>
          <div className="modal-body">
            {this.props.children}
          </div>
        </div>
      </div>
    </div>;
  }
}
